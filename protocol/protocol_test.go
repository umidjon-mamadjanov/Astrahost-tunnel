package protocol

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
)

func TestHeaderEncodeDecode(t *testing.T) {
	requestID := uuid.New()

	original := Header{
		Version:       Version,
		Type:          PacketConnect,
		Flags:         0x1234,
		RequestID:     requestID,
		PayloadLength: 1024,
	}

	data, err := original.Encode()
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	if len(data) != HeaderSize {
		t.Fatalf(
			"encoded header length = %d, want %d",
			len(data),
			HeaderSize,
		)
	}

	decoded, err := DecodeHeader(data)
	if err != nil {
		t.Fatalf("DecodeHeader() error = %v", err)
	}

	if decoded.Version != original.Version {
		t.Errorf(
			"Version = %d, want %d",
			decoded.Version,
			original.Version,
		)
	}

	if decoded.Type != original.Type {
		t.Errorf(
			"Type = %d, want %d",
			decoded.Type,
			original.Type,
		)
	}

	if decoded.Flags != original.Flags {
		t.Errorf(
			"Flags = %d, want %d",
			decoded.Flags,
			original.Flags,
		)
	}

	if decoded.RequestID != original.RequestID {
		t.Errorf(
			"RequestID = %s, want %s",
			decoded.RequestID,
			original.RequestID,
		)
	}

	if decoded.PayloadLength != original.PayloadLength {
		t.Errorf(
			"PayloadLength = %d, want %d",
			decoded.PayloadLength,
			original.PayloadLength,
		)
	}
}

func TestHeaderValidate(t *testing.T) {
	tests := []struct {
		name    string
		header  Header
		wantErr bool
	}{
		{
			name: "valid header",
			header: Header{
				Version:       Version,
				Type:          PacketConnect,
				RequestID:     uuid.New(),
				PayloadLength: 1024,
			},
			wantErr: false,
		},
		{
			name: "unsupported version",
			header: Header{
				Version:   99,
				Type:      PacketConnect,
				RequestID: uuid.New(),
			},
			wantErr: true,
		},
		{
			name: "payload too large",
			header: Header{
				Version:       Version,
				Type:          PacketConnect,
				RequestID:     uuid.New(),
				PayloadLength: MaxPayloadLength + 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.header.Validate()

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"Validate() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestDecodeHeaderRejectsShortData(t *testing.T) {
	data := make([]byte, HeaderSize-1)

	_, err := DecodeHeader(data)
	if err == nil {
		t.Fatal("DecodeHeader() expected error for short data")
	}
}

func TestEncodePacket(t *testing.T) {
	payload := []byte("hello astra")

	packet := Packet{
		Header: Header{
			Version:   Version,
			Type:      PacketConnect,
			RequestID: uuid.New(),
		},
		Payload: payload,
	}

	data, err := EncodePacket(packet)
	if err != nil {
		t.Fatalf("EncodePacket() error = %v", err)
	}

	if len(data) != HeaderSize+len(payload) {
		t.Fatalf(
			"encoded packet length = %d, want %d",
			len(data),
			HeaderSize+len(payload),
		)
	}

	header, err := DecodeHeader(data[:HeaderSize])
	if err != nil {
		t.Fatalf("DecodeHeader() error = %v", err)
	}

	if header.PayloadLength != uint32(len(payload)) {
		t.Fatalf(
			"PayloadLength = %d, want %d",
			header.PayloadLength,
			len(payload),
		)
	}

	if !bytes.Equal(data[HeaderSize:], payload) {
		t.Fatalf("payload mismatch")
	}
}

func TestEncodePacketRejectsOversizedPayload(t *testing.T) {
	payload := make([]byte, MaxPayloadLength+1)

	packet := Packet{
		Header: Header{
			Version:   Version,
			Type:      PacketConnect,
			RequestID: uuid.New(),
		},
		Payload: payload,
	}

	_, err := EncodePacket(packet)
	if err == nil {
		t.Fatal("EncodePacket() expected error for oversized payload")
	}
}

func TestConnectRequestRoundTrip(t *testing.T) {
	original := ConnectRequest{
		ProtocolVersion: Version,
		ClientVersion:   "1.0.0",
		Platform:        "linux",
		Architecture:    "amd64",
		TunnelName:      "demo",
		LocalHost:       "127.0.0.1",
		LocalPort:       5000,
		Auth:            "test-key",
	}

	data, err := EncodeConnectRequest(original)
	if err != nil {
		t.Fatalf("EncodeConnectRequest() error = %v", err)
	}

	decoded, err := DecodeConnectRequest(data)
	if err != nil {
		t.Fatalf("DecodeConnectRequest() error = %v", err)
	}

	if decoded != original {
		t.Fatalf(
			"decoded request = %+v, want %+v",
			decoded,
			original,
		)
	}
}

func TestConnectRequestValidation(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "empty payload",
			data: "",
		},
		{
			name: "invalid json",
			data: "{invalid",
		},
		{
			name: "unsupported version",
			data: `{"protocol_version":99,"client_version":"1.0.0","platform":"linux","architecture":"amd64","local_port":5000}`,
		},
		{
			name: "missing client version",
			data: `{"protocol_version":1,"platform":"linux","architecture":"amd64","local_port":5000}`,
		},
		{
			name: "missing platform",
			data: `{"protocol_version":1,"client_version":"1.0.0","architecture":"amd64","local_port":5000}`,
		},
		{
			name: "missing architecture",
			data: `{"protocol_version":1,"client_version":"1.0.0","platform":"linux","local_port":5000}`,
		},
		{
			name: "missing local port",
			data: `{"protocol_version":1,"client_version":"1.0.0","platform":"linux","architecture":"amd64"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := DecodeConnectRequest([]byte(tt.data)); err == nil {
				t.Fatal("DecodeConnectRequest() expected error")
			}
		})
	}
}

func TestConnectResponseRoundTrip(t *testing.T) {
	original := ConnectResponse{
		ProtocolVersion:   Version,
		SessionID:         "session-123",
		TunnelID:          "tunnel-123",
		HeartbeatInterval: 20,
		PublicURL:         "https://myapp.astra-tunnel.example",
	}

	data, err := EncodeConnectResponse(original)
	if err != nil {
		t.Fatalf("EncodeConnectResponse() error = %v", err)
	}

	decoded, err := DecodeConnectResponse(data)
	if err != nil {
		t.Fatalf("DecodeConnectResponse() error = %v", err)
	}

	if decoded != original {
		t.Fatalf(
			"decoded response = %+v, want %+v",
			decoded,
			original,
		)
	}
}

func TestConnectResponseValidation(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "empty payload",
			data: "",
		},
		{
			name: "invalid json",
			data: "{invalid",
		},
		{
			name: "unsupported version",
			data: `{"protocol_version":99,"session_id":"s","tunnel_id":"t","heartbeat_interval":20}`,
		},
		{
			name: "missing session id",
			data: `{"protocol_version":1,"tunnel_id":"t","heartbeat_interval":20}`,
		},
		{
			name: "missing tunnel id",
			data: `{"protocol_version":1,"session_id":"s","heartbeat_interval":20}`,
		},
		{
			name: "missing heartbeat interval",
			data: `{"protocol_version":1,"session_id":"s","tunnel_id":"t"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := DecodeConnectResponse([]byte(tt.data)); err == nil {
				t.Fatal("DecodeConnectResponse() expected error")
			}
		})
	}
}

func TestNewConnectPacket(t *testing.T) {
	req := ConnectRequest{
		ProtocolVersion: Version,
		ClientVersion:   "1.0.0",
		Platform:        "android",
		Architecture:    "arm64",
		TunnelName:      "demo",
		LocalHost:       "127.0.0.1",
		LocalPort:       5000,
		Auth:            "test-key",
	}

	packet, err := NewConnectPacket(req)
	if err != nil {
		t.Fatalf("NewConnectPacket() error = %v", err)
	}

	if packet.Header.Version != Version {
		t.Errorf(
			"Version = %d, want %d",
			packet.Header.Version,
			Version,
		)
	}

	if packet.Header.Type != PacketConnect {
		t.Errorf(
			"Type = %d, want %d",
			packet.Header.Type,
			PacketConnect,
		)
	}

	if packet.Header.RequestID == uuid.Nil {
		t.Fatal("RequestID must not be nil")
	}

	if int(packet.Header.PayloadLength) != len(packet.Payload) {
		t.Fatalf(
			"PayloadLength = %d, want %d",
			packet.Header.PayloadLength,
			len(packet.Payload),
		)
	}

	decoded, err := DecodeConnectRequest(packet.Payload)
	if err != nil {
		t.Fatalf("DecodeConnectRequest() error = %v", err)
	}

	if decoded != req {
		t.Fatalf(
			"decoded request = %+v, want %+v",
			decoded,
			req,
		)
	}
}

func TestNewConnectOKPacket(t *testing.T) {
	resp := ConnectResponse{
		ProtocolVersion:   Version,
		SessionID:         "session-123",
		TunnelID:          "tunnel-123",
		HeartbeatInterval: 20,
		PublicURL:         "https://myapp.astra-tunnel.example",
	}

	packet, err := NewConnectOKPacket(resp)
	if err != nil {
		t.Fatalf("NewConnectOKPacket() error = %v", err)
	}

	if packet.Header.Version != Version {
		t.Errorf(
			"Version = %d, want %d",
			packet.Header.Version,
			Version,
		)
	}

	if packet.Header.Type != PacketConnectOK {
		t.Errorf(
			"Type = %d, want %d",
			packet.Header.Type,
			PacketConnectOK,
		)
	}

	if packet.Header.RequestID == uuid.Nil {
		t.Fatal("RequestID must not be nil")
	}

	if int(packet.Header.PayloadLength) != len(packet.Payload) {
		t.Fatalf(
			"PayloadLength = %d, want %d",
			packet.Header.PayloadLength,
			len(packet.Payload),
		)
	}

	decoded, err := DecodeConnectResponse(packet.Payload)
	if err != nil {
		t.Fatalf("DecodeConnectResponse() error = %v", err)
	}

	if decoded != resp {
		t.Fatalf(
			"decoded response = %+v, want %+v",
			decoded,
			resp,
		)
	}
}

func TestPacketEncodeDecodeRoundTrip(t *testing.T) {
	req := ConnectRequest{
		ProtocolVersion: Version,
		ClientVersion:   "1.0.0",
		Platform:        "android",
		Architecture:    "arm64",
		TunnelName:      "demo",
		LocalHost:       "127.0.0.1",
		LocalPort:       5000,
		Auth:            "test-key",
	}

	packet, err := NewConnectPacket(req)
	if err != nil {
		t.Fatalf("NewConnectPacket() error = %v", err)
	}

	encoded, err := EncodePacket(packet)
	if err != nil {
		t.Fatalf("EncodePacket() error = %v", err)
	}

	decoded, err := DecodePacket(encoded)
	if err != nil {
		t.Fatalf("DecodePacket() error = %v", err)
	}

	if decoded.Header.Version != packet.Header.Version {
		t.Errorf(
			"Version = %d, want %d",
			decoded.Header.Version,
			packet.Header.Version,
		)
	}

	if decoded.Header.Type != packet.Header.Type {
		t.Errorf(
			"Type = %d, want %d",
			decoded.Header.Type,
			packet.Header.Type,
		)
	}

	if decoded.Header.Flags != packet.Header.Flags {
		t.Errorf(
			"Flags = %d, want %d",
			decoded.Header.Flags,
			packet.Header.Flags,
		)
	}

	if decoded.Header.RequestID != packet.Header.RequestID {
		t.Errorf(
			"RequestID = %s, want %s",
			decoded.Header.RequestID,
			packet.Header.RequestID,
		)
	}

	if decoded.Header.PayloadLength != packet.Header.PayloadLength {
		t.Errorf(
			"PayloadLength = %d, want %d",
			decoded.Header.PayloadLength,
			packet.Header.PayloadLength,
		)
	}

	if string(decoded.Payload) != string(packet.Payload) {
		t.Errorf("Payload mismatch")
	}
}

func TestDecodePacketRejectsInvalidLength(t *testing.T) {
	packet, err := NewConnectPacket(ConnectRequest{
		ProtocolVersion: Version,
		ClientVersion:   "1.0.0",
		Platform:        "linux",
		Architecture:    "amd64",
		LocalHost:       "127.0.0.1",
		LocalPort:       5000,
	})

	if err != nil {
		t.Fatalf("NewConnectPacket() error = %v", err)
	}

	encoded, err := EncodePacket(packet)
	if err != nil {
		t.Fatalf("EncodePacket() error = %v", err)
	}

	encoded = encoded[:len(encoded)-1]

	if _, err := DecodePacket(encoded); err == nil {
		t.Fatal("DecodePacket() expected error for invalid packet length")
	}
}

func TestPingPacket(t *testing.T) {
	packet := NewPingPacket()

	if packet.Header.Version != Version {
		t.Errorf(
			"Version = %d, want %d",
			packet.Header.Version,
			Version,
		)
	}

	if packet.Header.Type != PacketPing {
		t.Errorf(
			"Type = %d, want %d",
			packet.Header.Type,
			PacketPing,
		)
	}

	if packet.Header.RequestID == uuid.Nil {
		t.Fatal("PING RequestID must not be nil")
	}

	if err := ValidatePing(packet); err != nil {
		t.Fatalf("ValidatePing() error = %v", err)
	}

	encoded, err := EncodePacket(packet)
	if err != nil {
		t.Fatalf("EncodePacket() error = %v", err)
	}

	decoded, err := DecodePacket(encoded)
	if err != nil {
		t.Fatalf("DecodePacket() error = %v", err)
	}

	if err := ValidatePing(decoded); err != nil {
		t.Fatalf("ValidatePing(decoded) error = %v", err)
	}
}

func TestPongPacket(t *testing.T) {
	requestID := uuid.New()
	packet := NewPongPacket(requestID)

	if packet.Header.Version != Version {
		t.Errorf(
			"Version = %d, want %d",
			packet.Header.Version,
			Version,
		)
	}

	if packet.Header.Type != PacketPong {
		t.Errorf(
			"Type = %d, want %d",
			packet.Header.Type,
			PacketPong,
		)
	}

	if packet.Header.RequestID != requestID {
		t.Errorf(
			"RequestID = %s, want %s",
			packet.Header.RequestID,
			requestID,
		)
	}

	if err := ValidatePong(packet); err != nil {
		t.Fatalf("ValidatePong() error = %v", err)
	}

	encoded, err := EncodePacket(packet)
	if err != nil {
		t.Fatalf("EncodePacket() error = %v", err)
	}

	decoded, err := DecodePacket(encoded)
	if err != nil {
		t.Fatalf("DecodePacket() error = %v", err)
	}

	if err := ValidatePong(decoded); err != nil {
		t.Fatalf("ValidatePong(decoded) error = %v", err)
	}
}

func TestErrorPacket(t *testing.T) {
	requestID := uuid.New()

	packet, err := NewErrorPacket(
		requestID,
		ErrorAuthFailed,
		"authentication failed",
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := ValidateError(packet); err != nil {
		t.Fatal(err)
	}

	payload, err := DecodeError(packet.Payload)
	if err != nil {
		t.Fatal(err)
	}

	if payload.Code != ErrorAuthFailed {
		t.Fatalf(
			"unexpected error code: %s",
			payload.Code,
		)
	}

	if payload.Message != "authentication failed" {
		t.Fatalf(
			"unexpected error message: %s",
			payload.Message,
		)
	}

	if packet.Header.RequestID != requestID {
		t.Fatal("request ID mismatch")
	}
}

func TestClosePacket(t *testing.T) {
	requestID := uuid.New()

	packet, err := NewClosePacket(
		requestID,
		CloseClientShutdown,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := ValidateClose(packet); err != nil {
		t.Fatal(err)
	}

	payload, err := DecodeClose(packet.Payload)
	if err != nil {
		t.Fatal(err)
	}

	if payload.Reason != CloseClientShutdown {
		t.Fatalf(
			"unexpected close reason: %s",
			payload.Reason,
		)
	}

	if packet.Header.RequestID != requestID {
		t.Fatal("request ID mismatch")
	}
}

func TestHTTPRequestPacket(t *testing.T) {
	requestID := uuid.New()

	req := HTTPRequest{
		Method: "GET",
		URL:    "http://localhost:5000/test",
		Headers: []HTTPHeader{
			{
				Name:  "Host",
				Value: "localhost:5000",
			},
			{
				Name:  "User-Agent",
				Value: "Astra-Tunnel",
			},
		},
	}

	packet, err := NewHTTPRequestPacket(requestID, req)
	if err != nil {
		t.Fatal(err)
	}

	if err := ValidateHTTPRequest(packet); err != nil {
		t.Fatal(err)
	}

	decoded, err := DecodeHTTPRequest(packet.Payload)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.Method != req.Method {
		t.Fatalf(
			"method mismatch: got %s, want %s",
			decoded.Method,
			req.Method,
		)
	}

	if decoded.URL != req.URL {
		t.Fatalf(
			"URL mismatch: got %s, want %s",
			decoded.URL,
			req.URL,
		)
	}

	if packet.Header.RequestID != requestID {
		t.Fatal("request ID mismatch")
	}
}

func TestHTTPResponsePacket(t *testing.T) {
	requestID := uuid.New()

	resp := HTTPResponse{
		StatusCode: 200,
		Headers: []HTTPHeader{
			{
				Name:  "Content-Type",
				Value: "text/plain",
			},
		},
		Body: []byte("Hello from Astra Tunnel"),
	}

	packet, err := NewHTTPResponsePacket(requestID, resp)
	if err != nil {
		t.Fatal(err)
	}

	if err := ValidateHTTPResponse(packet); err != nil {
		t.Fatal(err)
	}

	decoded, err := DecodeHTTPResponse(packet.Payload)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.StatusCode != resp.StatusCode {
		t.Fatalf(
			"status code mismatch: got %d, want %d",
			decoded.StatusCode,
			resp.StatusCode,
		)
	}

	if string(decoded.Body) != string(resp.Body) {
		t.Fatalf(
			"body mismatch: got %q, want %q",
			decoded.Body,
			resp.Body,
		)
	}

	if packet.Header.RequestID != requestID {
		t.Fatal("request ID mismatch")
	}
}
