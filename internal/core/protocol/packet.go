package protocol

type Packet struct {
	Header  Header
	Payload []byte
}

func NewPacket(packetType PacketType, payload []byte) Packet {
	return Packet{
		Header: Header{
			Version: ProtocolVersion,
			Type:    packetType,
			Length:  uint32(len(payload)),
		},
		Payload: payload,
	}
}
