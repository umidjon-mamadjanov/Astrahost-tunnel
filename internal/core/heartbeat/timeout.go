package heartbeat

import (
	"time"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
)

const (
	HeartbeatInterval = 20 * time.Second
	HeartbeatTimeout  = 60 * time.Second
)

type TimeoutChecker struct {
	manager  *connection.Manager
	registry *connection.Registry
}

func NewTimeoutChecker(
	manager *connection.Manager,
	registry *connection.Registry,
) *TimeoutChecker {
	return &TimeoutChecker{
		manager:  manager,
		registry: registry,
	}
}

func (c *TimeoutChecker) Check(now time.Time) {
	for _, session := range c.manager.List() {
		if session.GetState() == connection.StateClosed {
			continue
		}

		if now.Sub(session.GetLastSeen()) <= HeartbeatTimeout {
			continue
		}

		session.SetState(connection.StateClosed)

		if session.Conn != nil {
			_ = session.Conn.Close()
		}

		tunnelID := session.GetTunnelID()

		if tunnelID != "" {
			c.registry.Unregister(tunnelID)
		}

		c.manager.Remove(session.ID)
	}
}
