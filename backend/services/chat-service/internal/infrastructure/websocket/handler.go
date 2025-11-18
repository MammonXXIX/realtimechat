package websocket

import (
	"log"
	"net/http"
	"realtimechat/services/chat-service/internal/domain"
	"realtimechat/shared/utils"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type HandlerWebsocket struct {
	Hub     *Hub
	Service domain.ChatService
}

func (h *HandlerWebsocket) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userid")
	if userID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Missing User ID",
			},
		})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket Upgrade Error: %v", err)
		return
	}

	client := &Client{
		hub:     h.Hub,
		conn:    conn,
		send:    make(chan *domain.MessageModel, 256),
		UserID:  userID,
		Rooms:   make(map[string]bool),
		Service: h.Service,
	}

	h.Hub.register <- client

	go client.ReadPump()
	go client.WritePump()
}
