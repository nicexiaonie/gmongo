package trace

import (
	"testing"

	"go.mongodb.org/mongo-driver/event"
)

func TestNewServerMonitor(t *testing.T) {
	called := map[string]bool{}
	monitor := NewServerMonitor(ServerMonitorOptions{
		ServerDescriptionChanged: func(event *event.ServerDescriptionChangedEvent) {
			called["serverDescriptionChanged"] = true
		},
		ServerOpening: func(event *event.ServerOpeningEvent) {
			called["serverOpening"] = true
		},
		ServerClosed: func(event *event.ServerClosedEvent) {
			called["serverClosed"] = true
		},
		TopologyDescriptionChanged: func(event *event.TopologyDescriptionChangedEvent) {
			called["topologyDescriptionChanged"] = true
		},
		TopologyOpening: func(event *event.TopologyOpeningEvent) {
			called["topologyOpening"] = true
		},
		TopologyClosed: func(event *event.TopologyClosedEvent) {
			called["topologyClosed"] = true
		},
		ServerHeartbeatStarted: func(event *event.ServerHeartbeatStartedEvent) {
			called["serverHeartbeatStarted"] = true
		},
		ServerHeartbeatSucceeded: func(event *event.ServerHeartbeatSucceededEvent) {
			called["serverHeartbeatSucceeded"] = true
		},
		ServerHeartbeatFailed: func(event *event.ServerHeartbeatFailedEvent) {
			called["serverHeartbeatFailed"] = true
		},
	})

	monitor.ServerDescriptionChanged(&event.ServerDescriptionChangedEvent{})
	monitor.ServerOpening(&event.ServerOpeningEvent{})
	monitor.ServerClosed(&event.ServerClosedEvent{})
	monitor.TopologyDescriptionChanged(&event.TopologyDescriptionChangedEvent{})
	monitor.TopologyOpening(&event.TopologyOpeningEvent{})
	monitor.TopologyClosed(&event.TopologyClosedEvent{})
	monitor.ServerHeartbeatStarted(&event.ServerHeartbeatStartedEvent{})
	monitor.ServerHeartbeatSucceeded(&event.ServerHeartbeatSucceededEvent{})
	monitor.ServerHeartbeatFailed(&event.ServerHeartbeatFailedEvent{})

	for name, ok := range called {
		if !ok {
			t.Fatalf("expected %s callback to be called", name)
		}
	}
	if len(called) != 9 {
		t.Fatalf("expected 9 callbacks to be called, got %d", len(called))
	}
}

func TestNewServerMonitorNilCallbacks(t *testing.T) {
	monitor := NewServerMonitor(ServerMonitorOptions{})

	monitor.ServerDescriptionChanged(&event.ServerDescriptionChangedEvent{})
	monitor.ServerOpening(&event.ServerOpeningEvent{})
	monitor.ServerClosed(&event.ServerClosedEvent{})
	monitor.TopologyDescriptionChanged(&event.TopologyDescriptionChangedEvent{})
	monitor.TopologyOpening(&event.TopologyOpeningEvent{})
	monitor.TopologyClosed(&event.TopologyClosedEvent{})
	monitor.ServerHeartbeatStarted(&event.ServerHeartbeatStartedEvent{})
	monitor.ServerHeartbeatSucceeded(&event.ServerHeartbeatSucceededEvent{})
	monitor.ServerHeartbeatFailed(&event.ServerHeartbeatFailedEvent{})
}
