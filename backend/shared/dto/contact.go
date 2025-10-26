package dto

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID        uuid.UUID `json:"id"`
	AdderID   string    `json:"adder_id"`
	AddedID   string    `json:"added_id"`
	AliasName string    `json:"alias_name"`
	CreatedAt time.Time `json:"created_at"`
}

type ContactWithAddedUserData struct {
	Contact
	AddedUser *AuthenticationData `json:"added_user"`
}
