package protocol

type ConnectRequest struct {
	ProtocolVersion uint8  `json:"protocol_version"`
	ClientVersion   string `json:"client_version"`
	Platform        string `json:"platform"`
	Architecture    string `json:"architecture"`
}

type ConnectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
