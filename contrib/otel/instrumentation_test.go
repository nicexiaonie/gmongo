package otel

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestInstrumentationCommandSpanSucceeded(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	instrumentation := New(Config{Enabled: true, TracerProvider: provider})
	monitor := instrumentation.CommandMonitor()

	ctx, parent := provider.Tracer("test").Start(context.Background(), "parent")
	command := rawCommand(t, bson.D{{Key: "find", Value: "users"}})

	monitor.Started(ctx, &event.CommandStartedEvent{
		Command:      command,
		DatabaseName: "myapp",
		CommandName:  "find",
		RequestID:    1,
		ConnectionID: "localhost:27017[-1]",
	})
	monitor.Succeeded(ctx, &event.CommandSucceededEvent{
		CommandFinishedEvent: event.CommandFinishedEvent{
			CommandName:  "find",
			DatabaseName: "myapp",
			RequestID:    1,
			ConnectionID: "localhost:27017[-1]",
		},
	})
	parent.End()

	spans := recorder.Ended()
	if len(spans) != 2 {
		t.Fatalf("expected 2 ended spans, got %d", len(spans))
	}

	mongoSpan := spans[0]
	if mongoSpan.Name() != "MongoDB find users" {
		t.Fatalf("unexpected span name: %s", mongoSpan.Name())
	}
	if mongoSpan.Parent().SpanID() != parent.SpanContext().SpanID() {
		t.Fatalf("expected mongo span to use parent context")
	}
	assertAttribute(t, mongoSpan.Attributes(), "db.system", "mongodb")
	assertAttribute(t, mongoSpan.Attributes(), "db.name", "myapp")
	assertAttribute(t, mongoSpan.Attributes(), "db.operation", "find")
	assertAttribute(t, mongoSpan.Attributes(), "db.mongodb.collection", "users")
}

func TestInstrumentationCommandSpanFailed(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	instrumentation := New(Config{Enabled: true, TracerProvider: provider})
	monitor := instrumentation.CommandMonitor()
	ctx := context.Background()

	monitor.Started(ctx, &event.CommandStartedEvent{
		Command:      rawCommand(t, bson.D{{Key: "insert", Value: "users"}}),
		DatabaseName: "myapp",
		CommandName:  "insert",
		RequestID:    2,
		ConnectionID: "localhost:27017[-1]",
	})
	monitor.Failed(ctx, &event.CommandFailedEvent{
		CommandFinishedEvent: event.CommandFinishedEvent{
			CommandName:  "insert",
			DatabaseName: "myapp",
			RequestID:    2,
			ConnectionID: "localhost:27017[-1]",
		},
		Failure: "duplicate key",
	})

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 ended span, got %d", len(spans))
	}
	if spans[0].Status().Code.String() != "Error" {
		t.Fatalf("expected error span status, got %s", spans[0].Status().Code.String())
	}
}

func TestInstrumentationUnmatchedFinishedEventsDoNotPanic(t *testing.T) {
	instrumentation := New(Config{Enabled: true, TracerProvider: otel.GetTracerProvider()})
	monitor := instrumentation.CommandMonitor()

	monitor.Succeeded(context.Background(), &event.CommandSucceededEvent{})
	monitor.Failed(context.Background(), &event.CommandFailedEvent{})
}

func rawCommand(t *testing.T, command bson.D) bson.Raw {
	t.Helper()
	raw, err := bson.Marshal(command)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func assertAttribute(t *testing.T, attrs []attribute.KeyValue, key string, value string) {
	t.Helper()
	for _, attr := range attrs {
		if string(attr.Key) == key && attr.Value.AsString() == value {
			return
		}
	}
	t.Fatalf("expected attribute %s=%s", key, value)
}
