package protocol

type Header struct {
	Version uint8
	Type    PacketType
	Length  uint32
}
