package handshake

import (
	"fmt"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/google/uuid"
)

type Handler struct {
	registry *connection.Registry
}

func New(registry *connection.Registry) *Handler {
	return &Handler{
		registry: registry,
	}
}

func (h *Handler) HandleConnect(
	session *connection.Session,
	req protocol.ConnectRequest,
) (*protocol.ConnectResponse, error) {
	if req.ProtocolVersion != protocol.Version {
		return nil, fmt.Errorf("unsupported protocol version")
	}

	session.ProtocolVersion = req.ProtocolVersion
	session.ClientVersion = req.ClientVersion
	session.Platform = req.Platform
	session.Architecture = req.Architecture
	session.State = connection.StateReady

	tunnelID := uuid.NewString()

	session.TunnelID = tunnelID

	h.registry.Register(tunnelID, session)

	return &protocol.ConnectResponse{
		ProtocolVersion:   protocol.Version,
		SessionID:         session.ID,
		TunnelID:          tunnelID,
		HeartbeatInterval: 20,
	}, nil
}
