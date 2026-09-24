package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/google/uuid"
)

func TestHTTPHandler(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf(
					"unexpected method: got %s, want GET",
					r.Method,
				)
			}

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Hello from localhost"))
		}),
	)
	defer server.Close()

	handler := NewHTTPHandler()

	requestID := uuid.New()

	request := protocol.HTTPRequest{
		Method: http.MethodGet,
		URL:    server.URL,
		Headers: []protocol.HTTPHeader{
			{
				Name:  "X-Tunnel-Test",
				Value: "astra",
			},
		},
	}

	payload, err := protocol.EncodeHTTPRequest(request)
	if err != nil {
		t.Fatalf("encode HTTP_REQUEST: %v", err)
	}

	packet := protocol.Packet{
		Header: protocol.Header{
			Version:   protocol.Version,
			Type:      protocol.PacketHTTPRequest,
			RequestID: requestID,
		},
		Payload: payload,
	}

	responsePacket, err := handler.Handle(packet)
	if err != nil {
		t.Fatalf("handle HTTP_REQUEST: %v", err)
	}

	if responsePacket.Header.Type != protocol.PacketHTTPResponse {
		t.Fatalf(
			"unexpected packet type: got %d, want %d",
			responsePacket.Header.Type,
			protocol.PacketHTTPResponse,
		)
	}

	if responsePacket.Header.RequestID != requestID {
		t.Fatalf(
			"request ID changed: got %s, want %s",
			responsePacket.Header.RequestID,
			requestID,
		)
	}

	response, err := protocol.DecodeHTTPResponse(
		responsePacket.Payload,
	)
	if err != nil {
		t.Fatalf("decode HTTP_RESPONSE: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf(
			"unexpected status: got %d, want %d",
			response.StatusCode,
			http.StatusOK,
		)
	}

	if string(response.Body) != "Hello from localhost" {
		t.Fatalf(
			"unexpected body: got %q",
			response.Body,
		)
	}
}
