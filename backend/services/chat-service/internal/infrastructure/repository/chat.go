package repository

import (
	"context"
	"database/sql"
	"errors"
	"realtimechat/services/chat-service/internal/domain"
	"realtimechat/shared/utils"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const queryTimeoutDuration = 5 * time.Second

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *chatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateChatRoom(ctx context.Context) (*domain.ChatRoomModel, error) {
	query := `
		INSERT INTO chat_rooms DEFAULT VALUES
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var chatRoom domain.ChatRoomModel

	err := r.db.QueryRowContext(ctx, query).Scan(
		&chatRoom.ID,
		&chatRoom.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &chatRoom, nil
}

func (r *chatRepository) CreateRoomMember(ctx context.Context, chatRoomID uuid.UUID, userIDs []string) error {
	query := `
		INSERT INTO chat_room_members (chat_room_id, user_id)
		VALUES ($1, $2), ($1, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		query,
		chatRoomID,
		userIDs[0],
		userIDs[1],
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return utils.ErrDuplicate
		}

		return err
	}

	return nil
}

func (r *chatRepository) CreateMessage(ctx context.Context, chatRoomID uuid.UUID, senderID, message string) error {
	query := `
		INSERT INTO messages (chat_room_id, sender_id, message)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		query,
		chatRoomID,
		senderID,
		message,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) GetChatHistoryByChatRoomID(ctx context.Context, chatRoomID uuid.UUID) ([]*domain.MessageModel, error) {
	query := `
		SELECT * FROM messages
		WHERE chat_room_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var messages []*domain.MessageModel

	rows, err := r.db.QueryContext(ctx, query, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := domain.MessageModel{}

		err := rows.Scan(
			&m.ID,
			&m.ChatRoomID,
			&m.SenderID,
			&m.Message,
			&m.IsRead,
			&m.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, &m)
	}

	return messages, nil
}

func (r *chatRepository) GetChatHistoriesByUserID(ctx context.Context, userID string) ([]*domain.ChatHistoryModel, error) {
	query := `
		SELECT cr.id AS chat_roo, other.user_id AS other_user_id, m.id, m.chat_room_id, m.sender_id, m.message, m.is_read, m.created_at
		FROM chat_rooms cr
		JOIN chat_room_members rm ON rm.chat_room_id = cr.id
		JOIN chat_room_members other ON other.chat_room_id = cr.id AND other.user_id != rm.user_id 
		LEFT JOIN LATERAL (
    	SELECT *
    	FROM messages
    	WHERE messages.chat_room_id = cr.id
    	ORDER BY created_at DESC
    	LIMIT 1
		) m ON TRUE WHERE rm.user_id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var chatHistories []*domain.ChatHistoryModel

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := domain.ChatHistoryModel{}
		m := domain.MessageModel{}

		var (
			messageID   sql.NullString
			chatRoomID  sql.NullString
			senderID    sql.NullString
			messageText sql.NullString
			isRead      sql.NullBool
			createdAt   sql.NullTime
		)

		if err := rows.Scan(
			&c.ChatRoomID,
			&c.OtherUserID,
			&messageID,
			&chatRoomID,
			&senderID,
			&messageText,
			&isRead,
			&createdAt,
		); err != nil {
			return nil, err
		}

		if messageID.Valid {
			m.ID, _ = uuid.Parse(messageID.String)
			m.ChatRoomID, _ = uuid.Parse(chatRoomID.String)
			m.SenderID = senderID.String
			m.Message = messageText.String
			m.IsRead = isRead.Bool
			m.CreatedAt = createdAt.Time

			c.LastMessage = &m
		} else {
			c.LastMessage = nil
		}

		chatHistories = append(chatHistories, &c)
	}

	return chatHistories, nil
}
