package heartbeat

import (
	"testing"
	"time"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
	"github.com/google/uuid"
)

func TestHandlePing(t *testing.T) {
	handler := New()

	oldLastSeen := time.Now().Add(-time.Minute)

	session := &connection.Session{
		ID:       "test-session",
		LastSeen: oldLastSeen,
	}

	requestID := uuid.New()

	packet := protocol.NewPingPacket()
	packet.Header.RequestID = requestID

	response, err := handler.HandlePing(session, &packet)
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("expected PONG response")
	}

	if response.Header.Type != protocol.PacketPong {
		t.Fatalf(
			"expected PONG packet, got %d",
			response.Header.Type,
		)
	}

	if response.Header.RequestID != requestID {
		t.Fatal("PONG request ID does not match PING request ID")
	}

	if !session.LastSeen.After(oldLastSeen) {
		t.Fatal("LastSeen was not updated")
	}

	if err := protocol.ValidatePong(*response); err != nil {
		t.Fatalf("invalid PONG packet: %v", err)
	}
}
