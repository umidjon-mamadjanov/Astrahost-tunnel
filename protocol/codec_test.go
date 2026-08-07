package protocol

import (
	"bytes"
	"testing"
)

func TestEncodeDecode(t *testing.T) {

	original := NewPacket(PacketPing, []byte("hello"))

	data, err := Encode(original)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.Header.Type != PacketPing {
		t.Fatal("packet type mismatch")
	}

	if !bytes.Equal(decoded.Payload, []byte("hello")) {
		t.Fatal("payload mismatch")
	}
}
