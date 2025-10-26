package dto

import "time"

type AuthenticationData struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}
