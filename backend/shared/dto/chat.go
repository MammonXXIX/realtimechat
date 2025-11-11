package dto

import (
	"time"

	"github.com/google/uuid"
)

type ChatRoom struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRoomMember struct {
	ID         uuid.UUID `json:"id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	UserID     string    `json:"user_id"`
	JoinedAt   time.Time `json:"joined_at"`
}

type Message struct {
	ID         uuid.UUID `json:"id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	SenderID   string    `json:"sender_id"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type ChatHistory struct {
	ChatRoomID   uuid.UUID `json:"chat_room_id"`
	ChatRoomName string    `json:"chat_room_name"`
	OtherUserID  string    `json:"other_user_id"`
	LastMessage  Message   `json:"last_message"`
}

type ChatHistoryWithOtherUserData struct {
	ChatHistory
	OtherUser *AuthenticationData `json:"other_user"`
}
