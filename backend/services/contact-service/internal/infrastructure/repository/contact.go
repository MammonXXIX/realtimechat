package repository

import (
	"context"
	"database/sql"
	"errors"
	"realtimechat/services/contact-service/internal/domain"
	"realtimechat/shared/utils"
	"time"

	"github.com/lib/pq"
)

const queryTimeoutDuration = 5 * time.Second

type contactRepository struct {
	db *sql.DB
}

func NewContactRepository(db *sql.DB) *contactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) CreateContact(ctx context.Context, contact *domain.ContactModel) error {
	query := `
		INSERT INTO contacts (adder_id, added_id, alias_name)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(
		ctx,
		query,
		contact.AdderID,
		contact.AddedID,
		contact.AliasName,
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

func (r *contactRepository) GetContacts(ctx context.Context, ID string) ([]*domain.ContactModel, error) {
	query := `
		SELECT * FROM contacts WHERE adder_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	var contacts []*domain.ContactModel

	rows, err := r.db.QueryContext(ctx, query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := domain.ContactModel{}

		err := rows.Scan(
			&c.ID,
			&c.AdderID,
			&c.AddedID,
			&c.AliasName,
			&c.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		contacts = append(contacts, &c)
	}

	return contacts, nil
}
