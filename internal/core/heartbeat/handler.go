package heartbeat

import (
	"fmt"
	"time"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) HandlePing(
	session *connection.Session,
	packet *protocol.Packet,
) (*protocol.Packet, error) {
	if err := protocol.ValidatePing(*packet); err != nil {
		return nil, fmt.Errorf("validate PING: %w", err)
	}

	session.LastSeen = time.Now()

	response := protocol.NewPongPacket(packet.Header.RequestID)

	return &response, nil
}

func (h *Handler) HandlePong(
	session *connection.Session,
	packet *protocol.Packet,
) (*protocol.Packet, error) {
	if err := protocol.ValidatePong(*packet); err != nil {
		return nil, fmt.Errorf("validate PONG: %w", err)
	}

	session.LastSeen = time.Now()

	return nil, nil
}
