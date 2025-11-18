package websocket

import (
	"log"
	"realtimechat/services/chat-service/internal/domain"
	"sync"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *domain.MessageModel
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *domain.MessageModel),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Printf("[Hub] Client %s Connected", client.UserID)
		case client := <-h.unregister:
			h.removeClient(client)
			log.Printf("[Hub] Client %s Disconnected", client.UserID)
		case message := <-h.broadcast:
			h.broadcastToRoom(message)
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]bool)
	}

	h.rooms[roomID][client] = true
	client.Rooms[roomID] = true

	log.Printf("[Hub] Client %s Joined Room %s", client.UserID, roomID)
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for roomID := range client.Rooms {
		if clients, ok := h.rooms[roomID]; ok {
			delete(clients, client)
		}
	}

	close(client.send)
}

func (h *Hub) broadcastToRoom(message *domain.MessageModel) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomID := message.ChatRoomID
	clients, ok := h.rooms[roomID.String()]

	if !ok {
		log.Printf("[Hub] Room %s Not Found", roomID)
	}

	log.Printf("[Hub] Broadcasting Message To Room %s (%d Clients)", roomID, len(clients))

	for client := range clients {
		select {
		case client.send <- message:
			log.Printf("[Hub] Message Sent To Client %s", client.UserID)
		default:
			log.Printf("[Hub] Client %s Buffer Full, Removing", client.UserID)
			h.unregister <- client
		}
	}
}

func (h *Hub) BroadcastMessage(message *domain.MessageModel) {
	h.broadcast <- message
}
