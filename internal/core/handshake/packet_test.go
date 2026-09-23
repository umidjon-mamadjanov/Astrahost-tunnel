package handshake

import (
	"testing"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

func TestHandleConnect(t *testing.T) {
	h := New()

	req := protocol.ConnectRequest{
		ProtocolVersion: protocol.Version,
		ClientVersion:   "1.0.0",
		Platform:        "android",
		Architecture:    "arm64",
		TunnelName:      "test",
	}

	packet, err := protocol.NewConnectPacket(req)
	if err != nil {
		t.Fatal(err)
	}

	session := &connection.Session{
		ID: "test-session",
	}

	err = h.HandlePacket(session, &packet)
	if err != nil {
		t.Fatal(err)
	}

	if session.ClientVersion != "1.0.0" {
		t.Fatalf(
			"client version not saved: got %s",
			session.ClientVersion,
		)
	}

	if session.ProtocolVersion != protocol.Version {
		t.Fatalf(
			"protocol version not saved: got %d",
			session.ProtocolVersion,
		)
	}

	if session.Platform != "android" {
		t.Fatalf(
			"platform not saved: got %s",
			session.Platform,
		)
	}

	if session.Architecture != "arm64" {
		t.Fatalf(
			"architecture not saved: got %s",
			session.Architecture,
		)
	}

	if session.State != connection.StateReady {
		t.Fatalf("session is not READY")
	}
}
