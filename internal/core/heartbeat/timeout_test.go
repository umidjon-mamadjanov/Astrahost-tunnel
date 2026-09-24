package heartbeat

import (
	"testing"
	"time"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
)

func TestTimeoutChecker(t *testing.T) {
	manager := connection.NewManager()
	registry := connection.NewRegistry()

	now := time.Now()

	expired := &connection.Session{
		ID:       "expired-session",
		State:    connection.StateReady,
		TunnelID: "expired-tunnel",
		LastSeen: now.Add(-61 * time.Second),
	}

	active := &connection.Session{
		ID:       "active-session",
		State:    connection.StateReady,
		TunnelID: "active-tunnel",
		LastSeen: now.Add(-59 * time.Second),
	}

	manager.Add(expired)
	manager.Add(active)

	registry.Register(expired.TunnelID, expired)
	registry.Register(active.TunnelID, active)

	checker := NewTimeoutChecker(manager, registry)
	checker.Check(now)

	if _, ok := manager.Get(expired.ID); ok {
		t.Fatal("expired session was not removed from manager")
	}

	if _, ok := registry.Get(expired.TunnelID); ok {
		t.Fatal("expired tunnel was not removed from registry")
	}

	if expired.State != connection.StateClosed {
		t.Fatal("expired session was not marked as closed")
	}

	if _, ok := manager.Get(active.ID); !ok {
		t.Fatal("active session was incorrectly removed")
	}

	if _, ok := registry.Get(active.TunnelID); !ok {
		t.Fatal("active tunnel was incorrectly removed")
	}
}
