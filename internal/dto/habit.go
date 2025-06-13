package dto

import "time"

type CreateHabitDTO struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Motivation  *string `json:"motivation" validate:"omitempty,max=1000"`
	Color       string  `json:"color" validate:"required,hexcolor"`
	Category    *string `json:"category" validate:"omitempty,max=100"`
	Frequency   string  `json:"frequency" validate:"required,oneof=daily weekly monthly custom"`
	TargetCount int     `json:"target_count" validate:"required,min=1"`
	TargetDays  *string `json:"target_days" validate:"omitempty,json"`
}

type HabitResponseDTO struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	Description      *string   `json:"description,omitempty"`
	Motivation       *string   `json:"motivation,omitempty"`
	Color            string    `json:"color"`
	Category         *string   `json:"category,omitempty"`
	Frequency        string    `json:"frequency"`
	TargetCount      int       `json:"target_count"`
	TargetDays       *string   `json:"target_days,omitempty"`
	CurrentStreak    int       `json:"current_streak"`
	BestStreak       int       `json:"best_streak"`
	TotalCompletions int       `json:"total_completions"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
