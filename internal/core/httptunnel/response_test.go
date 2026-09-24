package httptunnel

import (
	"testing"
	"time"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/google/uuid"
)

func TestResponseHandler(t *testing.T) {
	session := &connection.Session{
		ID:              "session-1",
		PendingRequests: connection.NewPendingRequests(),
	}

	handler := NewResponseHandler()

	requestID := uuid.New()
	response := protocol.HTTPResponse{
		StatusCode: 200,
		Headers: []protocol.HTTPHeader{
			{
				Name:  "Content-Type",
				Value: "text/plain",
			},
		},
		Body: []byte("Hello from localhost"),
	}

	payload, err := protocol.EncodeHTTPResponse(response)
	if err != nil {
		t.Fatalf("encode response: %v", err)
	}

	packet := protocol.Packet{
		Header: protocol.Header{
			Version:   protocol.Version,
			Type:      protocol.PacketHTTPResponse,
			RequestID: requestID,
		},
		Payload: payload,
	}

	ch := session.PendingRequests.Add(requestID)

	result, err := handler.Handle(session, &packet)
	if err != nil {
		t.Fatalf("handle HTTP_RESPONSE: %v", err)
	}

	if result != nil {
		t.Fatal("expected no response packet")
	}

	select {
	case data := <-ch:
		decoded, err := protocol.DecodeHTTPResponse(data)
		if err != nil {
			t.Fatalf("decode resolved response: %v", err)
		}

		if decoded.StatusCode != 200 {
			t.Fatalf(
				"unexpected status code: got %d, want 200",
				decoded.StatusCode,
			)
		}

		if string(decoded.Body) != "Hello from localhost" {
			t.Fatalf(
				"unexpected body: got %q",
				decoded.Body,
			)
		}

	case <-time.After(time.Second):
		t.Fatal("timed out waiting for HTTP response")
	}
}

func TestResponseHandlerUnknownRequest(t *testing.T) {
	session := &connection.Session{
		ID:              "session-1",
		PendingRequests: connection.NewPendingRequests(),
	}

	handler := NewResponseHandler()

	requestID := uuid.New()

	response := protocol.HTTPResponse{
		StatusCode: 200,
		Body:       []byte("test"),
	}

	payload, err := protocol.EncodeHTTPResponse(response)
	if err != nil {
		t.Fatalf("encode response: %v", err)
	}

	packet := protocol.Packet{
		Header: protocol.Header{
			Version:   protocol.Version,
			Type:      protocol.PacketHTTPResponse,
			RequestID: requestID,
		},
		Payload: payload,
	}

	_, err = handler.Handle(session, &packet)
	if err == nil {
		t.Fatal("expected error for unknown request ID")
	}
}
