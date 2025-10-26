package repository

import (
	"context"
	"database/sql"
	"realtimechat/services/authentication-service/internal/domain"
	"time"

	"github.com/lib/pq"
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
		INSERT INTO users (id, first_name, last_name, email, image_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.ImageURL,
		user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserModel, error) {
	query := `
		SELECT id, first_name, last_name, email, image_url, created_at FROM users WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var user domain.UserModel

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ImageURL,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUsersByIDs(ctx context.Context, IDs []string) ([]*domain.UserModel, error) {
	query := `
		SELECT id, first_name, last_name, email, image_url, created_at
		FROM users
		WHERE id = ANY($1)
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var users []*domain.UserModel

	rows, err := r.db.QueryContext(ctx, query, pq.Array(IDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.UserModel

		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.ImageURL,
			&u.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	return users, nil
}
