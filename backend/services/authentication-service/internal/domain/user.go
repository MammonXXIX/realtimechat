package domain

import (
	"context"
	"time"
)

type UserModel struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *UserModel) error
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	GetUsersByIDs(ctx context.Context, IDs []string) ([]*UserModel, error)
}

type UserService interface {
	CreateAccountByClerk(ctx context.Context, data any) error
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	GetUsersByIDs(ctx context.Context, IDs []string) ([]*UserModel, error)
}
