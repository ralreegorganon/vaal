package api

type joinMatchRequestMessage struct {
	Endpoint string `json:"endpoint"`
	Match    string `json:"match"`
}

type createMatchResponseMessage struct {
	Match string `json:"match"`
}

type startMatchRequestMessage struct {
	Match string `json:"match"`
}
