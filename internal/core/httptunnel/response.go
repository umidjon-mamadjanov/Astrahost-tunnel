package httptunnel

import (
	"fmt"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

type ResponseHandler struct{}

func NewResponseHandler() *ResponseHandler {
	return &ResponseHandler{}
}

func (h *ResponseHandler) Handle(
	session *connection.Session,
	packet *protocol.Packet,
) (*protocol.Packet, error) {
	if err := protocol.ValidateHTTPResponse(*packet); err != nil {
		return nil, fmt.Errorf("validate HTTP_RESPONSE: %w", err)
	}

	if session.PendingRequests == nil {
		return nil, fmt.Errorf("pending request manager is not initialized")
	}

	if !session.PendingRequests.Resolve(
		packet.Header.RequestID,
		packet.Payload,
	) {
		return nil, fmt.Errorf(
			"pending request not found: %s",
			packet.Header.RequestID,
		)
	}

	return nil, nil
}
