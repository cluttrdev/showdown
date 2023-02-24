package web

type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
