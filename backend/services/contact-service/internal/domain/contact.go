package domain

import (
	"context"
	"realtimechat/shared/dto"
	"time"

	"github.com/google/uuid"
)

type ContactModel struct {
	ID        uuid.UUID `json:"id"`
	AdderID   string    `json:"adder_id"`
	AddedID   string    `json:"added_id"`
	AliasName string    `json:"alias_name"`
	CreatedAt time.Time `json:"created_at"`
}

type ContactRepository interface {
	CreateContact(ctx context.Context, contact *ContactModel) error
	GetContacts(ctx context.Context, ID string) ([]*ContactModel, error)
}

type ContactService interface {
	CreateContactByEmail(ctx context.Context, adderID, email, aliasName string) error
	GetContactsByUserID(ctx context.Context, ID string) ([]*dto.ContactWithAddedUserData, error)
}
