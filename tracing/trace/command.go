package trace

import (
	"context"

	"go.mongodb.org/mongo-driver/event"
)

// CommandMonitorOptions 命令监控器配置
// 命令事件会携带执行 MongoDB 操作时传入的 context。
type CommandMonitorOptions struct {
	Started   func(context.Context, *event.CommandStartedEvent)
	Succeeded func(context.Context, *event.CommandSucceededEvent)
	Failed    func(context.Context, *event.CommandFailedEvent)
}

// NewCommandMonitor 创建命令监控器
func NewCommandMonitor(opts CommandMonitorOptions) *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(ctx context.Context, event *event.CommandStartedEvent) {
			if opts.Started != nil {
				opts.Started(ctx, event)
			}
		},
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			if opts.Succeeded != nil {
				opts.Succeeded(ctx, event)
			}
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			if opts.Failed != nil {
				opts.Failed(ctx, event)
			}
		},
	}
}
