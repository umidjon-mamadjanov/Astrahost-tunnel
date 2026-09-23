package handshake

import (
	"testing"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/protocol"
)

func TestHandleConnect(t *testing.T) {
	registry := connection.NewRegistry()
	handler := New(registry)

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

	response, err := handler.HandlePacket(session, &packet)
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("expected CONNECT_OK response")
	}

	if response.Header.Type != protocol.PacketConnectOK {
		t.Fatalf(
			"unexpected response type: %d",
			response.Header.Type,
		)
	}

	connectOK, err := protocol.DecodeConnectResponse(response.Payload)
	if err != nil {
		t.Fatal(err)
	}

	if connectOK.ProtocolVersion != protocol.Version {
		t.Fatalf("unexpected protocol version: %d", connectOK.ProtocolVersion)
	}

	if connectOK.SessionID != "test-session" {
		t.Fatalf("unexpected session ID: %s", connectOK.SessionID)
	}

	if connectOK.TunnelID == "" {
		t.Fatal("tunnel ID is empty")
	}

	registeredSession, ok := registry.Get(connectOK.TunnelID)
	if !ok {
		t.Fatal("tunnel was not registered")
	}

	if registeredSession != session {
		t.Fatal("registered session does not match")
	}

	if connectOK.HeartbeatInterval != 20 {
		t.Fatalf(
			"unexpected heartbeat interval: %d",
			connectOK.HeartbeatInterval,
		)
	}

	if response.Header.RequestID != packet.Header.RequestID {
		t.Fatal("request ID was not preserved")
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
		t.Fatal("session is not READY")
	}
}
