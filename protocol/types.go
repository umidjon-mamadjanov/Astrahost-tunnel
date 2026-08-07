package protocol

type PacketType uint8

const (
	PacketConnect PacketType = iota + 1
	PacketAuth
	PacketSubdomainCheck
	PacketSubdomainOK
	PacketSubdomainTaken
	PacketRegister
	PacketHTTPRequest
	PacketHTTPResponse
	PacketPing
	PacketPong
	PacketError
	PacketDisconnect
)
