package protocol

import "testing"

func TestConnectRequestLocalTarget(t *testing.T) {
	req := ConnectRequest{
		ProtocolVersion: Version,
		ClientVersion:   "1.0.0",
		Platform:        "linux",
		Architecture:    "arm64",
		LocalHost:       "127.0.0.1",
		LocalPort:       5000,
	}

	data, err := EncodeConnectRequest(req)
	if err != nil {
		t.Fatalf("encode CONNECT request: %v", err)
	}

	got, err := DecodeConnectRequest(data)
	if err != nil {
		t.Fatalf("decode CONNECT request: %v", err)
	}

	if got.LocalHost != "127.0.0.1" {
		t.Fatalf(
			"unexpected local host: got %q",
			got.LocalHost,
		)
	}

	if got.LocalPort != 5000 {
		t.Fatalf(
			"unexpected local port: got %d",
			got.LocalPort,
		)
	}
}

func TestConnectRequestDefaultLocalHost(t *testing.T) {
	req := ConnectRequest{
		ProtocolVersion: Version,
		ClientVersion:   "1.0.0",
		Platform:        "linux",
		Architecture:    "arm64",
		LocalPort:       3000,
	}

	data, err := EncodeConnectRequest(req)
	if err != nil {
		t.Fatalf("encode CONNECT request: %v", err)
	}

	got, err := DecodeConnectRequest(data)
	if err != nil {
		t.Fatalf("decode CONNECT request: %v", err)
	}

	if got.LocalHost != "127.0.0.1" {
		t.Fatalf(
			"expected default local host 127.0.0.1, got %q",
			got.LocalHost,
		)
	}

	if got.LocalPort != 3000 {
		t.Fatalf(
			"expected local port 3000, got %d",
			got.LocalPort,
		)
	}
}

func TestConnectRequestRequiresLocalPort(t *testing.T) {
	req := ConnectRequest{
		ProtocolVersion: Version,
		ClientVersion:   "1.0.0",
		Platform:        "linux",
		Architecture:    "arm64",
		LocalHost:       "127.0.0.1",
	}

	data, err := EncodeConnectRequest(req)
	if err != nil {
		t.Fatalf("encode CONNECT request: %v", err)
	}

	if _, err := DecodeConnectRequest(data); err == nil {
		t.Fatal("expected local_port validation error")
	}
}

func TestConnectResponsePublicURL(t *testing.T) {
	resp := ConnectResponse{
		ProtocolVersion:   Version,
		SessionID:         "session-1",
		TunnelID:          "tunnel-1",
		HeartbeatInterval: 20,
		PublicURL:         "https://myapp.astra-tunnel.example",
	}

	data, err := EncodeConnectResponse(resp)
	if err != nil {
		t.Fatalf("encode CONNECT_OK: %v", err)
	}

	got, err := DecodeConnectResponse(data)
	if err != nil {
		t.Fatalf("decode CONNECT_OK: %v", err)
	}

	if got.PublicURL != resp.PublicURL {
		t.Fatalf(
			"unexpected public URL: got %q, want %q",
			got.PublicURL,
			resp.PublicURL,
		)
	}
}
