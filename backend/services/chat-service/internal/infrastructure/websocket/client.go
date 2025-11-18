package websocket

import (
	"context"
	"encoding/json"
	"log"
	"realtimechat/services/chat-service/internal/domain"
	"time"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub     *Hub
	conn    *ws.Conn
	send    chan *domain.MessageModel
	UserID  string
	Rooms   map[string]bool
	Service domain.ChatService
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		var req map[string]any
		if err := json.Unmarshal(data, &req); err != nil {
			continue
		}

		switch req["type"] {
		case "join_rooms":
			rooms := req["rooms"].([]any)
			for _, r := range rooms {
				c.hub.JoinRoom(c, r.(string))
			}
		case "send_message":
			c.handleSendMessage(req)
		}
	}
}

func (c *Client) handleSendMessage(req map[string]any) {
	chatRoomIDStr, ok := req["room_id"].(string)
	if !ok {
		log.Printf("[Client] Missing room_id")
		return
	}

	chatRoomID, err := uuid.Parse(chatRoomIDStr)
	if err != nil {
		log.Printf("[Client] Invalid room_id format: %v", err)
		return
	}

	message, ok := req["message"].(string)
	if !ok {
		log.Printf("[Client] Missing messsage")
		return
	}

	reMessage, err := c.Service.CreateMessage(context.Background(), chatRoomID, c.UserID, message)
	if err != nil {
		log.Printf("[Client] Failed To Save: %v", err)
		return
	}

	c.hub.BroadcastMessage(reMessage)
}

func (c *Client) WritePump() {
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
				c.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Printf("[Client] JSON Marshal Error: %v", err)
				return
			}

			if err := c.conn.WriteMessage(ws.TextMessage, jsonData); err != nil {
				log.Printf("[Client] Write Error: %v", err)
				return
			}

			log.Printf("[Client] Message Delivered To %s", c.UserID)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
