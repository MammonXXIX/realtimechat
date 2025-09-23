package service

import (
	"context"
	"realtimechat/services/authentication-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type userService struct {
	repository domain.UserRepository
}

func NewUserService(repository domain.UserRepository) *userService {
	return &userService{repository: repository}
}

func (s *userService) Register(ctx context.Context, user *domain.UserModel) error {
	uuid, err := uuid.NewV6()
	if err != nil {
		return err
	}

	now := time.Now()

	user.UUID = uuid
	user.CreatedAt = now
	user.UpdatedAt = now

	return s.repository.CreateUser(ctx, user)
}
