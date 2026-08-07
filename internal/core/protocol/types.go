package protocol

type PacketType uint8

const (
	PacketConnect PacketType = iota + 1
	PacketConnectOK

	PacketAuth
	PacketAuthOK

	PacketPing
	PacketPong

	PacketHTTPRequest
	PacketHTTPResponse

	PacketError
	PacketDisconnect
)
