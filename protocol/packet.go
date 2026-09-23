package protocol

type Packet struct {
	Header  Header
	Payload []byte
}
