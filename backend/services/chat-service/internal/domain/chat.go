package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChatRoomModel struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRoomMemberModel struct {
	ID         uuid.UUID `json:"id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	UserID     string    `json:"user_id"`
	JoinedAt   time.Time `json:"joined_at"`
}

type MessageModel struct {
	ID         uuid.UUID `json:"id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	SenderID   string    `json:"sender_id"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type ChatHistoryModel struct {
	ChatRoomID   uuid.UUID `json:"chat_room_id"`
	ChatRoomName uuid.UUID `json:"chat_room_name"`
	OtherUserID  string    `json:"other_user_id"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
}

type ChatRepository interface {
	CreateChatRoom(ctx context.Context) (*ChatRoomModel, error)
	CreateRoomMember(ctx context.Context, chatRoomID uuid.UUID, userIDs []string) error
	CreateMessage(ctx context.Context, chatRoomID uuid.UUID, senderID, message string) error
	GetChatHistoryByChatRoomID(ctx context.Context, chatRoomID uuid.UUID) ([]*MessageModel, error)
	GetChatHistoriesByUserID(ctx context.Context, userID string) ([]*ChatHistoryModel, error)
}

type ChatService interface {
	CreatePrivateChat(ctx context.Context, userA, userB string) error
	CreateMessage(ctx context.Context, chatRoomID uuid.UUID, senderID, message string) error
	GetChatHistoryByChatRoomID(ctx context.Context, chatRoomID uuid.UUID) ([]*MessageModel, error)
	GetChatHistoriesByUserID(ctx context.Context, userID string) ([]*ChatHistoryModel, error)
}
