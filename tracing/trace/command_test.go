package trace

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/event"
)

func TestNewCommandMonitor(t *testing.T) {
	var started, succeeded, failed bool
	ctx := context.WithValue(context.Background(), struct{}{}, "trace")

	monitor := NewCommandMonitor(CommandMonitorOptions{
		Started: func(got context.Context, event *event.CommandStartedEvent) {
			if got != ctx {
				t.Fatalf("unexpected started context")
			}
			started = true
		},
		Succeeded: func(got context.Context, event *event.CommandSucceededEvent) {
			if got != ctx {
				t.Fatalf("unexpected succeeded context")
			}
			succeeded = true
		},
		Failed: func(got context.Context, event *event.CommandFailedEvent) {
			if got != ctx {
				t.Fatalf("unexpected failed context")
			}
			failed = true
		},
	})

	monitor.Started(ctx, &event.CommandStartedEvent{})
	monitor.Succeeded(ctx, &event.CommandSucceededEvent{})
	monitor.Failed(ctx, &event.CommandFailedEvent{})

	if !started || !succeeded || !failed {
		t.Fatalf("expected all command callbacks to be called")
	}
}

func TestNewCommandMonitorNilCallbacks(t *testing.T) {
	monitor := NewCommandMonitor(CommandMonitorOptions{})

	monitor.Started(context.Background(), &event.CommandStartedEvent{})
	monitor.Succeeded(context.Background(), &event.CommandSucceededEvent{})
	monitor.Failed(context.Background(), &event.CommandFailedEvent{})
}
