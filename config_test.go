package gmongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testMonitorProvider struct {
	commandMonitor *event.CommandMonitor
	poolMonitor    *event.PoolMonitor
	serverMonitor  *event.ServerMonitor
}

func (p testMonitorProvider) CommandMonitor() *event.CommandMonitor {
	return p.commandMonitor
}

func (p testMonitorProvider) PoolMonitor() *event.PoolMonitor {
	return p.poolMonitor
}

func (p testMonitorProvider) ServerMonitor() *event.ServerMonitor {
	return p.serverMonitor
}

func TestConfigToClientOptionsSetsMonitors(t *testing.T) {
	commandMonitor := &event.CommandMonitor{}
	poolMonitor := &event.PoolMonitor{}
	serverMonitor := &event.ServerMonitor{}

	config := DefaultConfig()
	config.CommandMonitor = commandMonitor
	config.PoolMonitor = poolMonitor
	config.ServerMonitor = serverMonitor

	opts := config.ToClientOptions()

	if opts.Monitor == nil {
		t.Fatalf("expected command monitor to be set")
	}
	if opts.PoolMonitor == nil {
		t.Fatalf("expected pool monitor to be set")
	}
	if opts.ServerMonitor == nil {
		t.Fatalf("expected server monitor to be set")
	}
}

func TestConfigToClientOptionsKeepsLegacyMonitor(t *testing.T) {
	called := false
	monitor := &event.CommandMonitor{
		Started: func(context.Context, *event.CommandStartedEvent) {
			called = true
		},
	}

	config := DefaultConfig()
	config.Monitor = monitor

	opts := config.ToClientOptions()
	opts.Monitor.Started(context.Background(), &event.CommandStartedEvent{})

	if !called {
		t.Fatalf("expected legacy monitor to be called")
	}
}

func TestConfigToClientOptionsCommandMonitorOverridesLegacyMonitor(t *testing.T) {
	legacyCalled := false
	commandCalled := false
	legacyMonitor := &event.CommandMonitor{
		Started: func(context.Context, *event.CommandStartedEvent) {
			legacyCalled = true
		},
	}
	commandMonitor := &event.CommandMonitor{
		Started: func(context.Context, *event.CommandStartedEvent) {
			commandCalled = true
		},
	}

	config := DefaultConfig()
	config.Monitor = legacyMonitor
	config.CommandMonitor = commandMonitor

	opts := config.ToClientOptions()
	opts.Monitor.Started(context.Background(), &event.CommandStartedEvent{})

	if legacyCalled {
		t.Fatalf("expected legacy monitor not to be called")
	}
	if !commandCalled {
		t.Fatalf("expected command monitor to be called")
	}
}

func TestConfigToClientOptionsComposesTracingAndManualMonitors(t *testing.T) {
	traceCalled := false
	manualCalled := false

	config := DefaultConfig()
	config.Tracing = testMonitorProvider{
		commandMonitor: &event.CommandMonitor{
			Started: func(context.Context, *event.CommandStartedEvent) {
				traceCalled = true
			},
		},
	}
	config.CommandMonitor = &event.CommandMonitor{
		Started: func(context.Context, *event.CommandStartedEvent) {
			manualCalled = true
		},
	}

	opts := config.ToClientOptions()
	opts.Monitor.Started(context.Background(), &event.CommandStartedEvent{})

	if !traceCalled || !manualCalled {
		t.Fatalf("expected tracing and manual monitors to be called")
	}
}

func TestConfigToClientOptionsHookRunsLast(t *testing.T) {
	commandMonitor := &event.CommandMonitor{}
	hookMonitor := &event.CommandMonitor{}

	config := DefaultConfig()
	config.CommandMonitor = commandMonitor
	config.ClientOptionsHook = func(opts *options.ClientOptions) {
		opts.SetMonitor(hookMonitor)
	}

	opts := config.ToClientOptions()

	if opts.Monitor != hookMonitor {
		t.Fatalf("expected client options hook to run last")
	}
}
