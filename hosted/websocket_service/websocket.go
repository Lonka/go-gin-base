package websocket_service

import (
	"go_gin_base/pkg/wsocket"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

const (
	namespace = "bbl"
)

var (
	hub = wsocket.NewHub(nil)
)

func Setup() {
	fn := func(message string, hub *wsocket.Hub) error {
		log.Println("message:", message)
		return nil
	}
	hub.OnMessage = fn
	go hub.Run()
}

func Serve(c *gin.Context) {
	wsocket.Serve(hub, c.Writer, c.Request)
}

func getWsMessage(event, message string) wsocket.WsMessage {
	return wsocket.WsMessage{
		Namespace: namespace,
		Event:     event,
		Message:   message,
	}
}

func heartbeat() {
	for {
		time.Sleep(10 * time.Second)
		message := getWsMessage(HEARTBEAT, "abc")
		hub.Broadcast <- message
	}
}

func GetClientLen() string {
	return com.ToStr(hub.Len)
}

func Broadcast(namespace, event string, message interface{}) {
	hub.Broadcast <- wsocket.WsMessage{
		Namespace: namespace,
		Event:     event,
		Message:   message,
	}
}
