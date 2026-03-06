package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

func (c *Client) SafeWriteJSON(v any) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.Conn.WriteJSON(v)
}

func (c *Client) SafeWriteMessage(messageType int, data []byte) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.Conn.WriteMessage(messageType, data)
}
