package handshake

import (
	"fmt"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/protocol"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) HandleConnect(
	session *connection.Session,
	req protocol.ConnectRequest,
) (*protocol.ConnectResponse, error) {

	if req.ProtocolVersion != protocol.ProtocolVersion {
		return nil, fmt.Errorf("unsupported protocol version")
	}

	session.State = connection.StateReady

	return &protocol.ConnectResponse{
		Success:       true,
		ServerVersion: protocol.ServerVersion,
		Message:       "Handshake completed",
	}, nil
}
