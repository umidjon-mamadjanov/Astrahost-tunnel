package protocol

import "testing"

func TestPacketEncodeDecode(t *testing.T) {

	payload := []byte("hello")

	packet := NewPacket(PacketConnect, payload)

	data, err := Encode(packet)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.Header.Type != PacketConnect {
		t.Fatal("packet type mismatch")
	}

	if string(decoded.Payload) != "hello" {
		t.Fatal("payload mismatch")
	}
}
