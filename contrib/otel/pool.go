package otel

import "go.mongodb.org/mongo-driver/event"

// PoolMonitorOptions 连接池监控器配置
// 连接池事件不关联单个 MongoDB 操作，因此不携带 context。
type PoolMonitorOptions struct {
	Event func(*event.PoolEvent)
}

// NewPoolMonitor 创建连接池监控器
func NewPoolMonitor(opts PoolMonitorOptions) *event.PoolMonitor {
	return &event.PoolMonitor{
		Event: func(event *event.PoolEvent) {
			if opts.Event != nil {
				opts.Event(event)
			}
		},
	}
}
