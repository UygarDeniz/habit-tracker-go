package entity

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	GoogleID  string    `json:"-"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(id, email, name, picture, googleID string) (*User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email cannot be empty")
	}
	if strings.TrimSpace(googleID) == "" {
		return nil, errors.New("googleID cannot be empty")
	}

	now := time.Now()
	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Picture:   picture,
		GoogleID:  googleID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
