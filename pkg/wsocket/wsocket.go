package wsocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	logging "go_gin_base/hosted/logging_service"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHub(onEvent OnMessageFunc) *Hub {
	return &Hub{
		Broadcast:  make(chan WsMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		OnMessage:  onEvent,
	}
}

type OnMessageFunc func(message string, hub *Hub) error

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	ip        string
	key       string
	namespace string
	events    map[string]byte
}

//Hub maintains the set of active clients and broadcasts messages to the
//clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	Broadcast chan WsMessage

	OnMessage OnMessageFunc

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	Len int
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			logging.WebSocket.Info(fmt.Sprintf("websocket broadcast message:%s.", message))
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		logging.WebSocket.Info(fmt.Sprintf("websocket ip:%s key:%s is closed.", c.ip, c.key))
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage() // wait client message
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logging.WebSocket.Error(fmt.Sprintf("websocket.IsUnexpectedCloseError:%v", err))
			}
			logging.WebSocket.Error("websocket read message have err")
			break
		}
		logging.WebSocket.Info(fmt.Sprintf("websocket ip:%s key:%s read message:%s.", c.ip, c.key, message))
		mssageStr := string(message)
		// 設定 client 的 namespace & event
		if strings.Contains(mssageStr, "initial-") {
			type initialSetting struct {
				Namespace string   `json:"namespace"`
				Events    []string `json:"events"`
			}
			var iSetting initialSetting
			err = json.Unmarshal([]byte(strings.Replace(mssageStr, "initial-", "", 1)), &iSetting)
			if err == nil {
				c.namespace = iSetting.Namespace
				if iSetting.Events != nil {
					c.events = make(map[string]byte)
					for _, v := range iSetting.Events {
						c.events[v] = 0x00
					}
				}
			}
		}
		if c.hub.OnMessage != nil {
			if err := c.hub.OnMessage(mssageStr, c.hub); err != nil {
				logging.WebSocket.Error(fmt.Sprintf("%v", err))
				break
			}
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.Len++
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.Len--
			}
		case message := <-h.Broadcast:
			sendMessage, _ := json.Marshal(message)
			for client := range h.clients {

				if client.namespace != "" && client.namespace != message.Namespace {
					logging.WebSocket.Info(fmt.Sprintf("message namespace: %s not equal client namespace: %s ", message.Namespace, client.namespace))
					continue
				}

				if client.events != nil {
					if _, ok := client.events[message.Event]; !ok {
						logging.WebSocket.Info(fmt.Sprintf("client events not contains message event: %s ", message.Event))

						continue
					}
				}
				select {
				case client.send <- sendMessage:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func Serve(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.WebSocket.Error(fmt.Sprintf("%v", err))
		return
	}
	socketKey := r.Header.Get("Sec-WebSocket-Key")
	ip := r.Header.Get("X-Real-IP")
	logging.WebSocket.Info(fmt.Sprintf("websocket ip:%s key:%s is connected.", ip, socketKey))
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), ip: ip, key: socketKey}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()

}
