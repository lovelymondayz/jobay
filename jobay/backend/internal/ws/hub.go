package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	slugClients map[string]map[*websocket.Conn]bool
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*websocket.Conn]bool),
		slugClients: make(map[string]map[*websocket.Conn]bool),
	}
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade error: %v", err)
		return
	}
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		for slug := range h.slugClients {
			delete(h.slugClients[slug], conn)
		}
		h.mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *Hub) HandleUserWS(w http.ResponseWriter, r *http.Request, slug string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade error: %v", err)
		return
	}
	h.mu.Lock()
	if h.slugClients[slug] == nil {
		h.slugClients[slug] = make(map[*websocket.Conn]bool)
	}
	h.slugClients[slug][conn] = true
	h.clients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		if h.slugClients[slug] != nil {
			delete(h.slugClients[slug], conn)
		}
		h.mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *Hub) Broadcast(msg interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		client.WriteJSON(msg)
	}
}

func (h *Hub) BroadcastToSlug(slug string, msg interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if clients, ok := h.slugClients[slug]; ok {
		for client := range clients {
			client.WriteJSON(msg)
		}
	}
}
