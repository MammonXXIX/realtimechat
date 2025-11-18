package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"realtimechat/services/chat-service/internal/domain"
	"realtimechat/shared/dto"
	"realtimechat/shared/utils"

	"github.com/google/uuid"
)

type chatService struct {
	repository domain.ChatRepository
}

func NewChatService(repository domain.ChatRepository) *chatService {
	return &chatService{repository: repository}
}

func (s *chatService) CreatePrivateChat(ctx context.Context, userA, userB string) error {
	chatRoom, err := s.repository.CreateChatRoom(ctx)
	if err != nil {
		return err
	}

	userIDs := []string{userA, userB}

	if err := s.repository.CreateRoomMember(ctx, chatRoom.ID, userIDs); err != nil {
		if errors.Is(err, utils.ErrDuplicate) {
			return utils.ErrDuplicate
		}

		return err

	}

	return nil
}

func (s *chatService) CreateMessage(ctx context.Context, chatRoomID uuid.UUID, senderID, message string) (*domain.MessageModel, error) {
	reMessage, err := s.repository.CreateMessage(ctx, chatRoomID, senderID, message)
	if err != nil {
		return nil, err
	}

	return reMessage, nil
}

func (s *chatService) GetChatHistoryByChatRoomID(ctx context.Context, chatRoomID uuid.UUID) ([]*domain.MessageModel, error) {
	messages, err := s.repository.GetChatHistoryByChatRoomID(ctx, chatRoomID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		return nil, err
	}

	return messages, err
}

func (s *chatService) GetChatHistoriesByUserID(ctx context.Context, userID string) ([]*dto.ChatHistoryWithOtherUserData, error) {
	chatHistories, err := s.repository.GetChatHistoriesByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		return nil, err
	}

	userIDs := make([]string, 0, len(chatHistories))
	for _, c := range chatHistories {
		userIDs = append(userIDs, c.OtherUserID)
	}

	requestBody, err := json.Marshal(map[string][]string{"users_ids": userIDs})
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequestWithContext(ctx, "POST", "http://authentication-service:8082/bulk", bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resBody []dto.AuthenticationData
	if err := utils.DecodeJSON(res.Body, &resBody); err != nil {
		return nil, err
	}

	userMap := make(map[string]dto.AuthenticationData)
	for _, u := range resBody {
		userMap[u.ID] = u
	}

	result := make([]*dto.ChatHistoryWithOtherUserData, 0, len(chatHistories))
	for _, c := range chatHistories {
		var lastMessage dto.Message

		if c.LastMessage != nil {
			lastMessage = dto.Message{
				ID:         c.LastMessage.ID,
				ChatRoomID: c.LastMessage.ChatRoomID,
				SenderID:   c.LastMessage.SenderID,
				Message:    c.LastMessage.Message,
				IsRead:     c.LastMessage.IsRead,
				CreatedAt:  c.LastMessage.CreatedAt,
			}
		}

		chatHistory := &dto.ChatHistoryWithOtherUserData{
			ChatHistory: dto.ChatHistory{
				ChatRoomID:  c.ChatRoomID,
				OtherUserID: c.OtherUserID,
				LastMessage: lastMessage,
			},
		}

		if user, ok := userMap[c.OtherUserID]; ok {
			chatHistory.ChatRoomName = user.FirstName
			chatHistory.OtherUser = &user
		}

		result = append(result, chatHistory)
	}

	return result, err
}
