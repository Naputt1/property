package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"backend/internal/config"
	"backend/internal/repository"
	"backend/internal/routes/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Manager struct {
	clients    map[*Client]struct{}
	broadcast  chan interface{}
	register   chan *Client
	unregister chan *Client
	done       chan struct{}
	mutex      sync.Mutex
	cfg        *config.Config
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

func NewManager(cfg *config.Config) repository.SocketService {
	return &Manager{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		done:       make(chan struct{}),
		cfg:        cfg,
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client] = struct{}{}
			m.mutex.Unlock()
		case client := <-m.unregister:
			m.mutex.Lock()
			delete(m.clients, client)
			m.mutex.Unlock()
			client.Conn.Close()
		case msg := <-m.broadcast:
			m.mutex.Lock()
			for client := range m.clients {
				err := client.SafeWriteJSON(msg)
				if err != nil {
					client.Conn.Close()
					delete(m.clients, client)
				}
			}
			m.mutex.Unlock()
		case <-m.done:
			return
		}
	}
}

func (m *Manager) Broadcast(message interface{}) {
	m.broadcast <- message
}

func (m *Manager) ServeWS(c *gin.Context) {
	_, exists := c.Get(config.CONTEXT_USER)
	if !exists {
		middlewares.ReturnUnauth(c, m.cfg)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		Conn: conn,
	}

	m.register <- client

	// Reader loop (app-level ping/pong)
	go func() {
		defer func() {
			m.unregister <- client
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WS read error: %v", err)
				}
				break
			}

			// Handle app-level ping/pong
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				if msg["type"] == "ping" {
					_ = client.SafeWriteJSON(map[string]string{"type": "pong"})
					continue
				}
			}
		}
	}()

	// Protocol-level pinger
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := client.SafeWriteMessage(websocket.PingMessage, nil); err != nil {
					m.mutex.Lock()
					delete(m.clients, client)
					m.mutex.Unlock()
					client.Conn.Close()
					return
				}
			case <-m.done:
				return
			}
		}
	}()
}
