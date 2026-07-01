package otel

import "go.mongodb.org/mongo-driver/event"

// ServerMonitorOptions 服务发现和心跳监控器配置
// 服务、拓扑、心跳事件不关联单个 MongoDB 操作，因此不携带 context。
type ServerMonitorOptions struct {
	ServerDescriptionChanged   func(*event.ServerDescriptionChangedEvent)
	ServerOpening              func(*event.ServerOpeningEvent)
	ServerClosed               func(*event.ServerClosedEvent)
	TopologyDescriptionChanged func(*event.TopologyDescriptionChangedEvent)
	TopologyOpening            func(*event.TopologyOpeningEvent)
	TopologyClosed             func(*event.TopologyClosedEvent)
	ServerHeartbeatStarted     func(*event.ServerHeartbeatStartedEvent)
	ServerHeartbeatSucceeded   func(*event.ServerHeartbeatSucceededEvent)
	ServerHeartbeatFailed      func(*event.ServerHeartbeatFailedEvent)
}

// NewServerMonitor 创建服务发现和心跳监控器
func NewServerMonitor(opts ServerMonitorOptions) *event.ServerMonitor {
	return &event.ServerMonitor{
		ServerDescriptionChanged: func(event *event.ServerDescriptionChangedEvent) {
			if opts.ServerDescriptionChanged != nil {
				opts.ServerDescriptionChanged(event)
			}
		},
		ServerOpening: func(event *event.ServerOpeningEvent) {
			if opts.ServerOpening != nil {
				opts.ServerOpening(event)
			}
		},
		ServerClosed: func(event *event.ServerClosedEvent) {
			if opts.ServerClosed != nil {
				opts.ServerClosed(event)
			}
		},
		TopologyDescriptionChanged: func(event *event.TopologyDescriptionChangedEvent) {
			if opts.TopologyDescriptionChanged != nil {
				opts.TopologyDescriptionChanged(event)
			}
		},
		TopologyOpening: func(event *event.TopologyOpeningEvent) {
			if opts.TopologyOpening != nil {
				opts.TopologyOpening(event)
			}
		},
		TopologyClosed: func(event *event.TopologyClosedEvent) {
			if opts.TopologyClosed != nil {
				opts.TopologyClosed(event)
			}
		},
		ServerHeartbeatStarted: func(event *event.ServerHeartbeatStartedEvent) {
			if opts.ServerHeartbeatStarted != nil {
				opts.ServerHeartbeatStarted(event)
			}
		},
		ServerHeartbeatSucceeded: func(event *event.ServerHeartbeatSucceededEvent) {
			if opts.ServerHeartbeatSucceeded != nil {
				opts.ServerHeartbeatSucceeded(event)
			}
		},
		ServerHeartbeatFailed: func(event *event.ServerHeartbeatFailedEvent) {
			if opts.ServerHeartbeatFailed != nil {
				opts.ServerHeartbeatFailed(event)
			}
		},
	}
}
