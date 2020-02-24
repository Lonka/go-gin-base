package wsocket

type WsMessage struct {
	Namespace string      `json:"namespace"`
	Event     string      `json:"event"`
	Message   interface{} `json:"message"`
}
