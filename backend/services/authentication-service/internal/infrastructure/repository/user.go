package repository

import (
	"context"
	"database/sql"
	"realtimechat/services/authentication-service/internal/domain"
	"time"
)

const queryTimeoutDuration = 5 * time.Second

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.UserModel) error {
	query := `
		INSERT INTO users (uuid, first_name, last_name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.UUID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// user.LastName = "repository"

	return nil
}
