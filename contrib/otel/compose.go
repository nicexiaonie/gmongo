package otel

import (
	"context"

	"go.mongodb.org/mongo-driver/event"
)

// ComposeCommandMonitors 组合多个命令监控器
func ComposeCommandMonitors(monitors ...*event.CommandMonitor) *event.CommandMonitor {
	monitors = compactCommandMonitors(monitors)
	if len(monitors) == 0 {
		return nil
	}
	return &event.CommandMonitor{
		Started: func(ctx context.Context, event *event.CommandStartedEvent) {
			for _, monitor := range monitors {
				if monitor.Started != nil {
					monitor.Started(ctx, event)
				}
			}
		},
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			for _, monitor := range monitors {
				if monitor.Succeeded != nil {
					monitor.Succeeded(ctx, event)
				}
			}
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			for _, monitor := range monitors {
				if monitor.Failed != nil {
					monitor.Failed(ctx, event)
				}
			}
		},
	}
}

// ComposePoolMonitors 组合多个连接池监控器
func ComposePoolMonitors(monitors ...*event.PoolMonitor) *event.PoolMonitor {
	monitors = compactPoolMonitors(monitors)
	if len(monitors) == 0 {
		return nil
	}
	return &event.PoolMonitor{
		Event: func(event *event.PoolEvent) {
			for _, monitor := range monitors {
				if monitor.Event != nil {
					monitor.Event(event)
				}
			}
		},
	}
}

// ComposeServerMonitors 组合多个服务发现和心跳监控器
func ComposeServerMonitors(monitors ...*event.ServerMonitor) *event.ServerMonitor {
	monitors = compactServerMonitors(monitors)
	if len(monitors) == 0 {
		return nil
	}
	return &event.ServerMonitor{
		ServerDescriptionChanged: func(event *event.ServerDescriptionChangedEvent) {
			for _, monitor := range monitors {
				if monitor.ServerDescriptionChanged != nil {
					monitor.ServerDescriptionChanged(event)
				}
			}
		},
		ServerOpening: func(event *event.ServerOpeningEvent) {
			for _, monitor := range monitors {
				if monitor.ServerOpening != nil {
					monitor.ServerOpening(event)
				}
			}
		},
		ServerClosed: func(event *event.ServerClosedEvent) {
			for _, monitor := range monitors {
				if monitor.ServerClosed != nil {
					monitor.ServerClosed(event)
				}
			}
		},
		TopologyDescriptionChanged: func(event *event.TopologyDescriptionChangedEvent) {
			for _, monitor := range monitors {
				if monitor.TopologyDescriptionChanged != nil {
					monitor.TopologyDescriptionChanged(event)
				}
			}
		},
		TopologyOpening: func(event *event.TopologyOpeningEvent) {
			for _, monitor := range monitors {
				if monitor.TopologyOpening != nil {
					monitor.TopologyOpening(event)
				}
			}
		},
		TopologyClosed: func(event *event.TopologyClosedEvent) {
			for _, monitor := range monitors {
				if monitor.TopologyClosed != nil {
					monitor.TopologyClosed(event)
				}
			}
		},
		ServerHeartbeatStarted: func(event *event.ServerHeartbeatStartedEvent) {
			for _, monitor := range monitors {
				if monitor.ServerHeartbeatStarted != nil {
					monitor.ServerHeartbeatStarted(event)
				}
			}
		},
		ServerHeartbeatSucceeded: func(event *event.ServerHeartbeatSucceededEvent) {
			for _, monitor := range monitors {
				if monitor.ServerHeartbeatSucceeded != nil {
					monitor.ServerHeartbeatSucceeded(event)
				}
			}
		},
		ServerHeartbeatFailed: func(event *event.ServerHeartbeatFailedEvent) {
			for _, monitor := range monitors {
				if monitor.ServerHeartbeatFailed != nil {
					monitor.ServerHeartbeatFailed(event)
				}
			}
		},
	}
}

func compactCommandMonitors(monitors []*event.CommandMonitor) []*event.CommandMonitor {
	result := make([]*event.CommandMonitor, 0, len(monitors))
	for _, monitor := range monitors {
		if monitor != nil {
			result = append(result, monitor)
		}
	}
	return result
}

func compactPoolMonitors(monitors []*event.PoolMonitor) []*event.PoolMonitor {
	result := make([]*event.PoolMonitor, 0, len(monitors))
	for _, monitor := range monitors {
		if monitor != nil {
			result = append(result, monitor)
		}
	}
	return result
}

func compactServerMonitors(monitors []*event.ServerMonitor) []*event.ServerMonitor {
	result := make([]*event.ServerMonitor, 0, len(monitors))
	for _, monitor := range monitors {
		if monitor != nil {
			result = append(result, monitor)
		}
	}
	return result
}
