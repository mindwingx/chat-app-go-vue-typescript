package entity

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn     *websocket.Conn
	username string
	send     chan []byte
	once     sync.Once
}

func NewWebSocketClient(conn *websocket.Conn, username string) *WebSocketClient {
	return &WebSocketClient{
		conn:     conn,
		username: username,
		send:     make(chan []byte, 100),
	}
}

func (c *WebSocketClient) Conn() *websocket.Conn {
	return c.conn
}

func (c *WebSocketClient) Username() string {
	return c.username
}

func (c *WebSocketClient) SendChannel() chan []byte {
	return c.send
}

func (c *WebSocketClient) Close() {
	c.once.Do(func() {
		close(c.send)
	})
}
