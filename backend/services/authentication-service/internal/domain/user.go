package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	UUID      uuid.UUID `json:"uuid"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *UserModel) error
	// GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	// UpdateUser(ctx context.Context, user *UserModel) error
	// DeleteUser(ctx context.Context, uuid uuid.UUID) error
}

type UserService interface {
	Register(ctx context.Context, user *UserModel) error
	// Login(ctx context.Context, email, password string) (*UserModel, error)
}
