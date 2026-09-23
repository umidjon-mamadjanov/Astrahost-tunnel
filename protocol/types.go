package protocol

const (
	Version uint8 = 1
)

type PacketType uint8

const (
	PacketConnect   PacketType = 0x01
	PacketConnectOK PacketType = 0x02

	PacketHTTPRequest  PacketType = 0x10
	PacketHTTPResponse PacketType = 0x11

	PacketPing PacketType = 0x20
	PacketPong PacketType = 0x21

	PacketError PacketType = 0x30
	PacketClose PacketType = 0x31
)

const (
	HeaderSize       = 24
	MaxPayloadLength = 16 * 1024 * 1024
)
