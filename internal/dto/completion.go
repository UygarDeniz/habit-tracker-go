package dto

import (
	"time"
)

// CreateCompletionDTO represents the request to create a habit completion
type CreateCompletionDTO struct {
	CompletionDate string  `json:"completion_date" validate:"required,datetime=2006-01-02"`
	Count          int     `json:"count" validate:"required,min=1"`
	Notes          *string `json:"notes" validate:"omitempty,max=1000"`
}

// UpdateCompletionDTO represents the request to update a habit completion
type UpdateCompletionDTO struct {
	Count *int    `json:"count,omitempty" validate:"omitempty,min=1"`
	Notes *string `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// CompletionResponseDTO represents the response containing habit completion details
type CompletionResponseDTO struct {
	ID             string    `json:"id"`
	HabitID        string    `json:"habit_id"`
	UserID         string    `json:"user_id"`
	CompletedAt    time.Time `json:"completed_at"`
	CompletionDate time.Time `json:"completion_date"`
	Count          int       `json:"count"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// GetCompletionsQueryDTO represents query parameters for getting completions
type GetCompletionsQueryDTO struct {
	HabitID   *string `json:"habit_id"`
	StartDate *string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   *string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	Limit     *int    `json:"limit" validate:"omitempty,min=1,max=1000"`
	Offset    *int    `json:"offset" validate:"omitempty,min=0"`
}
