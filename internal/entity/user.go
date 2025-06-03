package entity

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUser(id, email, username, passwordHash string) (*User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email cannot be empty")
	}

	now := time.Now()
	return &User{
		ID:           id,
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
