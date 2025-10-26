package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"realtimechat/services/authentication-service/internal/domain"
	"realtimechat/shared/utils"
	"time"
)

type userService struct {
	repository domain.UserRepository
}

func NewUserService(repository domain.UserRepository) *userService {
	return &userService{repository: repository}
}

func (s *userService) CreateAccountByClerk(ctx context.Context, data any) error {
	d, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("Invalid Clerk Data Format")
	}

	id, _ := d["id"].(string)
	firstName, _ := d["first_name"].(string)

	lastName := ""
	if i, ok := d["external_accounts"].([]any); ok && len(i) > 0 {
		if j, ok := i[0].(map[string]any); ok {
			lastName, _ = j["family_name"].(string)
		}
	}

	email := ""
	if i, ok := d["email_addresses"].([]any); ok && len(i) > 0 {
		if j, ok := i[0].(map[string]any); ok {
			email, _ = j["email_address"].(string)
		}
	}

	imageUrl, _ := d["image_url"].(string)
	createdAt := time.Now()

	user := &domain.UserModel{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		ImageURL:  imageUrl,
		CreatedAt: createdAt,
	}

	return s.repository.CreateUser(ctx, user)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.UserModel, error) {
	user, err := s.repository.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *userService) GetUsersByIDs(ctx context.Context, IDs []string) ([]*domain.UserModel, error) {
	users, err := s.repository.GetUsersByIDs(ctx, IDs)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		return nil, err
	}

	return users, nil
}
