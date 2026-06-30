package trace

import (
	"testing"

	"go.mongodb.org/mongo-driver/event"
)

func TestNewPoolMonitor(t *testing.T) {
	called := false
	monitor := NewPoolMonitor(PoolMonitorOptions{
		Event: func(event *event.PoolEvent) {
			called = true
		},
	})

	monitor.Event(&event.PoolEvent{})

	if !called {
		t.Fatalf("expected pool callback to be called")
	}
}

func TestNewPoolMonitorNilCallback(t *testing.T) {
	monitor := NewPoolMonitor(PoolMonitorOptions{})
	monitor.Event(&event.PoolEvent{})
}
