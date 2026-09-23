package protocol

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type ErrorCode string

const (
	ErrorProtocolVersionUnsupported ErrorCode = "PROTOCOL_VERSION_UNSUPPORTED"
	ErrorInvalidPacket              ErrorCode = "INVALID_PACKET"
	ErrorInvalidPayload             ErrorCode = "INVALID_PAYLOAD"
	ErrorPayloadTooLarge            ErrorCode = "PAYLOAD_TOO_LARGE"
	ErrorAuthRequired               ErrorCode = "AUTH_REQUIRED"
	ErrorAuthFailed                 ErrorCode = "AUTH_FAILED"
	ErrorTunnelNotFound             ErrorCode = "TUNNEL_NOT_FOUND"
	ErrorSessionNotFound            ErrorCode = "SESSION_NOT_FOUND"
	ErrorRequestTimeout             ErrorCode = "REQUEST_TIMEOUT"
	ErrorRateLimited                ErrorCode = "RATE_LIMITED"
	ErrorInternal                   ErrorCode = "INTERNAL_ERROR"
)

type ErrorPayload struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewErrorPacket(requestID uuid.UUID, code ErrorCode, message string) (Packet, error) {
	if code == "" {
		return Packet{}, fmt.Errorf("error code is required")
	}

	payload := ErrorPayload{
		Code:    code,
		Message: message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return Packet{}, fmt.Errorf("encode ERROR payload: %w", err)
	}

	return Packet{
		Header: Header{
			Version:       Version,
			Type:          PacketError,
			RequestID:     requestID,
			PayloadLength: uint32(len(data)),
		},
		Payload: data,
	}, nil
}

func DecodeError(data []byte) (ErrorPayload, error) {
	var payload ErrorPayload

	if len(data) == 0 {
		return payload, fmt.Errorf("empty ERROR payload")
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return payload, fmt.Errorf("decode ERROR payload: %w", err)
	}

	if payload.Code == "" {
		return payload, fmt.Errorf("error code is required")
	}

	return payload, nil
}

func ValidateError(packet Packet) error {
	if packet.Header.Type != PacketError {
		return fmt.Errorf("expected ERROR packet, got %d", packet.Header.Type)
	}

	if err := packet.Header.Validate(); err != nil {
		return err
	}

	_, err := DecodeError(packet.Payload)
	return err
}

type CloseReason string

const (
	CloseClientShutdown     CloseReason = "client_shutdown"
	CloseServerShutdown     CloseReason = "server_shutdown"
	CloseAuthenticationFail CloseReason = "authentication_failed"
	CloseProtocolError      CloseReason = "protocol_error"
	CloseTimeout            CloseReason = "timeout"
	CloseAdministrative     CloseReason = "administrative"
)

type ClosePayload struct {
	Reason CloseReason `json:"reason"`
}

func NewClosePacket(requestID uuid.UUID, reason CloseReason) (Packet, error) {
	if reason == "" {
		return Packet{}, fmt.Errorf("close reason is required")
	}

	payload := ClosePayload{
		Reason: reason,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return Packet{}, fmt.Errorf("encode CLOSE payload: %w", err)
	}

	return Packet{
		Header: Header{
			Version:       Version,
			Type:          PacketClose,
			RequestID:     requestID,
			PayloadLength: uint32(len(data)),
		},
		Payload: data,
	}, nil
}

func DecodeClose(data []byte) (ClosePayload, error) {
	var payload ClosePayload

	if len(data) == 0 {
		return payload, fmt.Errorf("empty CLOSE payload")
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return payload, fmt.Errorf("decode CLOSE payload: %w", err)
	}

	if payload.Reason == "" {
		return payload, fmt.Errorf("close reason is required")
	}

	return payload, nil
}

func ValidateClose(packet Packet) error {
	if packet.Header.Type != PacketClose {
		return fmt.Errorf("expected CLOSE packet, got %d", packet.Header.Type)
	}

	if err := packet.Header.Validate(); err != nil {
		return err
	}

	_, err := DecodeClose(packet.Payload)
	return err
}
