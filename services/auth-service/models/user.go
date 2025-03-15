package models

import (
	"github.com/google/uuid"
)

type User struct{
	ID uuid.UUID `json:"id"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt string `json:"created_at"`
}