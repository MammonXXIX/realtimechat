package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"realtimechat/services/chat-service/internal/domain"
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

func (s *chatService) CreateMessage(ctx context.Context, chatRoomID uuid.UUID, senderID, message string) error {
	if err := s.repository.CreateMessage(ctx, chatRoomID, senderID, message); err != nil {
		return err
	}

	return nil
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

func (s *chatService) GetChatHistoriesByUserID(ctx context.Context, userID string) ([]*domain.ChatHistoryModel, error) {
	chatHistories, err := s.repository.GetChatHistoriesByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		return nil, err
	}

	log.Println(chatHistories)

	return chatHistories, err
}
