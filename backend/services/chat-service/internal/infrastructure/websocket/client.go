package websocket

import (
	ws "github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	conn   *ws.Conn
	send   chan []byte
	RoomID string
	UserID string
}
