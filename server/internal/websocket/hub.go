package websocket

import (
	"sync"
	"thingify/server/internal/model"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) RegisterClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = true
}

func (h *Hub) UnRegisterClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

func (h *Hub) BroadcastIssue(issue model.Issue) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// TODO: отправка не каждому а определенному клиенту
	for conn := range h.clients {
		conn.WriteJSON(issue)
	}
}
