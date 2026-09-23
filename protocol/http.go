package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HTTPRequest struct {
	Method  string       `json:"method"`
	URL     string       `json:"url"`
	Headers []HTTPHeader `json:"headers,omitempty"`
	Body    []byte       `json:"body,omitempty"`
}

type HTTPResponse struct {
	StatusCode int          `json:"status_code"`
	Headers    []HTTPHeader `json:"headers,omitempty"`
	Body       []byte       `json:"body,omitempty"`
}

func EncodeHTTPRequest(req HTTPRequest) ([]byte, error) {
	if req.Method == "" {
		return nil, fmt.Errorf("HTTP method is required")
	}

	if req.URL == "" {
		return nil, fmt.Errorf("HTTP URL is required")
	}

	return json.Marshal(req)
}

func DecodeHTTPRequest(data []byte) (HTTPRequest, error) {
	var req HTTPRequest

	if len(data) == 0 {
		return req, fmt.Errorf("empty HTTP_REQUEST payload")
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return req, fmt.Errorf("decode HTTP_REQUEST payload: %w", err)
	}

	if req.Method == "" {
		return req, fmt.Errorf("HTTP method is required")
	}

	if req.URL == "" {
		return req, fmt.Errorf("HTTP URL is required")
	}

	return req, nil
}

func EncodeHTTPResponse(resp HTTPResponse) ([]byte, error) {
	if resp.StatusCode < 100 || resp.StatusCode > 599 {
		return nil, fmt.Errorf("invalid HTTP status code: %d", resp.StatusCode)
	}

	return json.Marshal(resp)
}

func DecodeHTTPResponse(data []byte) (HTTPResponse, error) {
	var resp HTTPResponse

	if len(data) == 0 {
		return resp, fmt.Errorf("empty HTTP_RESPONSE payload")
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, fmt.Errorf("decode HTTP_RESPONSE payload: %w", err)
	}

	if resp.StatusCode < 100 || resp.StatusCode > 599 {
		return resp, fmt.Errorf("invalid HTTP status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func NewHTTPRequestPacket(requestID uuid.UUID, req HTTPRequest) (Packet, error) {
	payload, err := EncodeHTTPRequest(req)
	if err != nil {
		return Packet{}, fmt.Errorf("encode HTTP_REQUEST: %w", err)
	}

	return Packet{
		Header: Header{
			Version:       Version,
			Type:          PacketHTTPRequest,
			RequestID:     requestID,
			PayloadLength: uint32(len(payload)),
		},
		Payload: payload,
	}, nil
}

func NewHTTPResponsePacket(requestID uuid.UUID, resp HTTPResponse) (Packet, error) {
	payload, err := EncodeHTTPResponse(resp)
	if err != nil {
		return Packet{}, fmt.Errorf("encode HTTP_RESPONSE: %w", err)
	}

	return Packet{
		Header: Header{
			Version:       Version,
			Type:          PacketHTTPResponse,
			RequestID:     requestID,
			PayloadLength: uint32(len(payload)),
		},
		Payload: payload,
	}, nil
}

func ValidateHTTPRequest(packet Packet) error {
	if packet.Header.Type != PacketHTTPRequest {
		return fmt.Errorf(
			"expected HTTP_REQUEST packet, got %d",
			packet.Header.Type,
		)
	}

	if err := packet.Header.Validate(); err != nil {
		return err
	}

	_, err := DecodeHTTPRequest(packet.Payload)
	return err
}

func ValidateHTTPResponse(packet Packet) error {
	if packet.Header.Type != PacketHTTPResponse {
		return fmt.Errorf(
			"expected HTTP_RESPONSE packet, got %d",
			packet.Header.Type,
		)
	}

	if err := packet.Header.Validate(); err != nil {
		return err
	}

	_, err := DecodeHTTPResponse(packet.Payload)
	return err
}
