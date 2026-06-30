package trace

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/event"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type commandKey struct {
	requestID    int64
	connectionID string
}

type commandSpan struct {
	span oteltrace.Span
}

// Instrumentation MongoDB trace instrumentation
type Instrumentation struct {
	config Config
	tracer oteltrace.Tracer

	mu    sync.Mutex
	spans map[commandKey]commandSpan
}

// New 创建 MongoDB trace instrumentation
func New(config Config) *Instrumentation {
	if config.TracerName == "" {
		config.TracerName = defaultTracerName
	}
	if config.TracerProvider == nil {
		config.TracerProvider = otel.GetTracerProvider()
	}

	return &Instrumentation{
		config: config,
		tracer: config.TracerProvider.Tracer(config.TracerName),
		spans:  make(map[commandKey]commandSpan),
	}
}

// CommandMonitor 创建会记录 OpenTelemetry span 的命令监控器
func (i *Instrumentation) CommandMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started:   i.commandStarted,
		Succeeded: i.commandSucceeded,
		Failed:    i.commandFailed,
	}
}

// PoolMonitor 创建连接池诊断监控器
func (i *Instrumentation) PoolMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{
		Event: func(event *event.PoolEvent) {
			if i.config.PoolEventHandler != nil {
				i.config.PoolEventHandler(PoolEventInfo{
					Type:         event.Type,
					Address:      event.Address,
					ConnectionID: event.ConnectionID,
					Reason:       event.Reason,
					Error:        event.Error,
				})
			}
		},
	}
}

// ServerMonitor 创建服务发现和心跳诊断监控器
func (i *Instrumentation) ServerMonitor() *event.ServerMonitor {
	return &event.ServerMonitor{
		ServerDescriptionChanged: func(event *event.ServerDescriptionChangedEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "server.description_changed"})
		},
		ServerOpening: func(event *event.ServerOpeningEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "server.opening"})
		},
		ServerClosed: func(event *event.ServerClosedEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "server.closed"})
		},
		TopologyDescriptionChanged: func(event *event.TopologyDescriptionChangedEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "topology.description_changed"})
		},
		TopologyOpening: func(event *event.TopologyOpeningEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "topology.opening"})
		},
		TopologyClosed: func(event *event.TopologyClosedEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "topology.closed"})
		},
		ServerHeartbeatStarted: func(event *event.ServerHeartbeatStartedEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "server.heartbeat.started", ConnectionID: event.ConnectionID})
		},
		ServerHeartbeatSucceeded: func(event *event.ServerHeartbeatSucceededEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "server.heartbeat.succeeded", ConnectionID: event.ConnectionID})
		},
		ServerHeartbeatFailed: func(event *event.ServerHeartbeatFailedEvent) {
			i.handleServerEvent(ServerEventInfo{Type: "server.heartbeat.failed", ConnectionID: event.ConnectionID, Failure: event.Failure})
		},
	}
}

func (i *Instrumentation) commandStarted(ctx context.Context, event *event.CommandStartedEvent) {
	if !i.config.Enabled {
		return
	}

	info := commandStartedInfo(event)
	_, span := i.tracer.Start(ctx, i.spanName(info), oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	span.SetAttributes(i.attributes(info)...)

	i.mu.Lock()
	i.spans[key(info.RequestID, info.ConnectionID)] = commandSpan{span: span}
	i.mu.Unlock()
}

func (i *Instrumentation) commandSucceeded(ctx context.Context, event *event.CommandSucceededEvent) {
	span, ok := i.takeSpan(event.RequestID, event.ConnectionID)
	if !ok {
		return
	}
	span.SetAttributes(attribute.Int64("db.mongodb.duration_ms", event.Duration.Milliseconds()))
	span.SetStatus(codes.Ok, "")
	span.End()
}

func (i *Instrumentation) commandFailed(ctx context.Context, event *event.CommandFailedEvent) {
	span, ok := i.takeSpan(event.RequestID, event.ConnectionID)
	if !ok {
		return
	}
	span.SetAttributes(attribute.Int64("db.mongodb.duration_ms", event.Duration.Milliseconds()))
	span.RecordError(errors.New(event.Failure))
	span.SetStatus(codes.Error, event.Failure)
	span.End()
}

func (i *Instrumentation) takeSpan(requestID int64, connectionID string) (oteltrace.Span, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	key := key(requestID, connectionID)
	span, ok := i.spans[key]
	if !ok {
		return nil, false
	}
	delete(i.spans, key)
	return span.span, true
}

func (i *Instrumentation) handleServerEvent(info ServerEventInfo) {
	if i.config.ServerEventHandler != nil {
		i.config.ServerEventHandler(info)
	}
}

func (i *Instrumentation) spanName(info CommandStartedInfo) string {
	if i.config.SpanNameFormatter != nil {
		return i.config.SpanNameFormatter(info)
	}
	if info.Collection != "" {
		return fmt.Sprintf("MongoDB %s %s", info.CommandName, info.Collection)
	}
	return fmt.Sprintf("MongoDB %s", info.CommandName)
}

func (i *Instrumentation) attributes(info CommandStartedInfo) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("db.system", defaultDBSystem),
		attribute.String("db.name", info.DatabaseName),
		attribute.String("db.operation", info.CommandName),
		attribute.Int64("db.mongodb.request_id", info.RequestID),
		attribute.String("db.mongodb.connection_id", info.ConnectionID),
	}
	if info.Collection != "" {
		attrs = append(attrs, attribute.String("db.mongodb.collection", info.Collection))
	}
	if i.config.Attributes != nil {
		attrs = append(attrs, i.config.Attributes(info)...)
	}
	if i.config.RecordCommand && i.config.CommandSanitizer != nil {
		attrs = append(attrs, i.config.CommandSanitizer(info)...)
	}
	return attrs
}

func commandStartedInfo(event *event.CommandStartedEvent) CommandStartedInfo {
	return CommandStartedInfo{
		DatabaseName: event.DatabaseName,
		CommandName:  event.CommandName,
		Collection:   collectionName(event),
		RequestID:    event.RequestID,
		ConnectionID: event.ConnectionID,
	}
}

func collectionName(event *event.CommandStartedEvent) string {
	value := event.Command.Lookup(event.CommandName)
	if collection, ok := value.StringValueOK(); ok {
		return collection
	}
	return ""
}

func key(requestID int64, connectionID string) commandKey {
	return commandKey{requestID: requestID, connectionID: connectionID}
}
