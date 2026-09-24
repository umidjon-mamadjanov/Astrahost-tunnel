package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type ConnectRequest struct {
	ProtocolVersion uint8  `json:"protocol_version"`
	ClientVersion   string `json:"client_version"`
	Platform        string `json:"platform"`
	Architecture    string `json:"architecture"`
	TunnelName      string `json:"tunnel_name,omitempty"`
	LocalHost       string `json:"local_host,omitempty"`
	LocalPort       uint16 `json:"local_port,omitempty"`
	Auth            string `json:"auth,omitempty"`
}

type ConnectResponse struct {
	ProtocolVersion   uint8  `json:"protocol_version"`
	SessionID         string `json:"session_id"`
	TunnelID          string `json:"tunnel_id"`
	HeartbeatInterval uint32 `json:"heartbeat_interval"`
	PublicURL         string `json:"public_url,omitempty"`
}

func EncodeConnectRequest(req ConnectRequest) ([]byte, error) {
	if req.ProtocolVersion == 0 {
		req.ProtocolVersion = Version
	}

	return json.Marshal(req)
}

func DecodeConnectRequest(data []byte) (ConnectRequest, error) {
	var req ConnectRequest

	if len(data) == 0 {
		return req, fmt.Errorf("empty CONNECT payload")
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return req, fmt.Errorf("decode CONNECT payload: %w", err)
	}

	if req.ProtocolVersion != Version {
		return req, fmt.Errorf(
			"unsupported protocol version: %d",
			req.ProtocolVersion,
		)
	}

	if req.ClientVersion == "" {
		return req, fmt.Errorf("client_version is required")
	}

	if req.Platform == "" {
		return req, fmt.Errorf("platform is required")
	}

	if req.Architecture == "" {
		return req, fmt.Errorf("architecture is required")
	}

	if req.LocalHost == "" {
		req.LocalHost = "127.0.0.1"
	}

	if req.LocalPort == 0 {
		return req, fmt.Errorf("local_port is required")
	}

	return req, nil
}

func EncodeConnectResponse(resp ConnectResponse) ([]byte, error) {
	if resp.ProtocolVersion == 0 {
		resp.ProtocolVersion = Version
	}

	return json.Marshal(resp)
}

func DecodeConnectResponse(data []byte) (ConnectResponse, error) {
	var resp ConnectResponse

	if len(data) == 0 {
		return resp, fmt.Errorf("empty CONNECT_OK payload")
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, fmt.Errorf("decode CONNECT_OK payload: %w", err)
	}

	if resp.ProtocolVersion != Version {
		return resp, fmt.Errorf(
			"unsupported protocol version: %d",
			resp.ProtocolVersion,
		)
	}

	if resp.SessionID == "" {
		return resp, fmt.Errorf("session_id is required")
	}

	if resp.TunnelID == "" {
		return resp, fmt.Errorf("tunnel_id is required")
	}

	if resp.HeartbeatInterval == 0 {
		return resp, fmt.Errorf("heartbeat_interval is required")
	}

	return resp, nil
}

func NewConnectPacket(req ConnectRequest) (Packet, error) {
	payload, err := EncodeConnectRequest(req)
	if err != nil {
		return Packet{}, fmt.Errorf("encode CONNECT: %w", err)
	}

	return Packet{
		Header: Header{
			Version:       Version,
			Type:          PacketConnect,
			RequestID:     uuid.New(),
			PayloadLength: uint32(len(payload)),
		},
		Payload: payload,
	}, nil
}

func NewConnectOKPacket(resp ConnectResponse) (Packet, error) {
	payload, err := EncodeConnectResponse(resp)
	if err != nil {
		return Packet{}, fmt.Errorf("encode CONNECT_OK: %w", err)
	}

	return Packet{
		Header: Header{
			Version:       Version,
			Type:          PacketConnectOK,
			RequestID:     uuid.New(),
			PayloadLength: uint32(len(payload)),
		},
		Payload: payload,
	}, nil
}
