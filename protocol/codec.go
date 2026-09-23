package protocol

import "fmt"

func EncodePacket(packet Packet) ([]byte, error) {
	payloadLength := len(packet.Payload)

	if payloadLength > MaxPayloadLength {
		return nil, fmt.Errorf(
			"payload too large: %d bytes, maximum %d",
			payloadLength,
			MaxPayloadLength,
		)
	}

	packet.Header.PayloadLength = uint32(payloadLength)

	header, err := packet.Header.Encode()
	if err != nil {
		return nil, err
	}

	result := make([]byte, 0, HeaderSize+payloadLength)
	result = append(result, header...)
	result = append(result, packet.Payload...)

	return result, nil
}

func DecodePacket(data []byte) (Packet, error) {
	if len(data) < HeaderSize {
		return Packet{}, fmt.Errorf(
			"packet too short: got %d bytes, minimum %d",
			len(data),
			HeaderSize,
		)
	}

	header, err := DecodeHeader(data[:HeaderSize])
	if err != nil {
		return Packet{}, fmt.Errorf("decode packet header: %w", err)
	}

	if err := header.Validate(); err != nil {
		return Packet{}, err
	}

	expectedLength := HeaderSize + int(header.PayloadLength)

	if len(data) != expectedLength {
		return Packet{}, fmt.Errorf(
			"invalid packet length: got %d, want %d",
			len(data),
			expectedLength,
		)
	}

	payload := make([]byte, header.PayloadLength)
	copy(payload, data[HeaderSize:])

	return Packet{
		Header:  header,
		Payload: payload,
	}, nil
}
