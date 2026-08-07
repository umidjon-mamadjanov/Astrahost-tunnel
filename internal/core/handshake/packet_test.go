package handshake

import (
	"testing"

	"github.com/astrahost/astrahost-tunnel/internal/core/connection"
	"github.com/astrahost/astrahost-tunnel/internal/core/protocol"
)

func TestHandleConnect(t *testing.T) {

	h := New()

	req := protocol.ConnectRequest{
		ProtocolVersion: protocol.ProtocolVersion,
		ClientVersion:   "0.1.0",
		Platform:        "android",
		Architecture:    "arm64",
	}

	payload, err := protocol.EncodeJSON(req)
	if err != nil {
		t.Fatal(err)
	}

	packet := protocol.NewPacket(
		protocol.PacketConnect,
		payload,
	)

	session := &connection.Session{}

	err = h.HandlePacket(session, &packet)
	if err != nil {
		t.Fatal(err)
	}

	if session.ClientVersion != "0.1.0" {
		t.Fatal("client version not saved")
	}
}
