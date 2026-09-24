package client

import (
	"fmt"

	"github.com/astrahost/astrahost-tunnel/protocol"
)

type HTTPHandler struct {
	forwarder *HTTPForwarder
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		forwarder: NewHTTPForwarder(),
	}
}

func (h *HTTPHandler) Handle(
	packet protocol.Packet,
) (protocol.Packet, error) {
	if packet.Header.Type != protocol.PacketHTTPRequest {
		return protocol.Packet{}, fmt.Errorf(
			"expected HTTP_REQUEST packet, got %d",
			packet.Header.Type,
		)
	}

	request, err := protocol.DecodeHTTPRequest(packet.Payload)
	if err != nil {
		return protocol.Packet{}, fmt.Errorf(
			"decode HTTP_REQUEST: %w",
			err,
		)
	}

	response, err := h.forwarder.Forward(request)
	if err != nil {
		return protocol.Packet{}, fmt.Errorf(
			"forward HTTP request: %w",
			err,
		)
	}

	responsePacket, err := protocol.NewHTTPResponsePacket(
		packet.Header.RequestID,
		response,
	)
	if err != nil {
		return protocol.Packet{}, fmt.Errorf(
			"create HTTP_RESPONSE: %w",
			err,
		)
	}

	return responsePacket, nil
}
