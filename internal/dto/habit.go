package dto

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/uygardeniz/habit-tracker/internal/entity"
)

type CreateHabitDTO struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Motivation  *string `json:"motivation" validate:"omitempty,max=1000"`
	Color       string  `json:"color" validate:"required,hexcolor"`
	Category    *string `json:"category" validate:"omitempty,max=100"`
	Frequency   string  `json:"frequency" validate:"required,oneof=daily weekly monthly"`
	TargetCount int     `json:"target_count" validate:"required,min=1"`
	TargetDays  *string `json:"target_days" validate:"omitempty,json"`
}

type UpdateHabitDTO struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Motivation  *string `json:"motivation,omitempty" validate:"omitempty,max=1000"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Category    *string `json:"category,omitempty" validate:"omitempty,max=100"`
	Frequency   *string `json:"frequency,omitempty" validate:"omitempty,oneof=daily weekly monthly"`
	TargetCount *int    `json:"target_count,omitempty" validate:"omitempty,min=1"`
	TargetDays  *string `json:"target_days,omitempty" validate:"omitempty,json"`
	IsActive    *bool   `json:"is_active,omitempty"`
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

// ConvertTargetDaysFromJSON converts JSON string to TargetDays entity
func ConvertTargetDaysFromJSON(targetDaysJSON *string) (*entity.TargetDays, error) {
	if targetDaysJSON == nil || *targetDaysJSON == "" {
		return nil, nil
	}

	var targetDaysData struct {
		Days []any `json:"days"`
	}

	if err := json.Unmarshal([]byte(*targetDaysJSON), &targetDaysData); err != nil {
		return nil, errors.New("invalid target days JSON format")
	}

	return &entity.TargetDays{Days: targetDaysData.Days}, nil
}

// ConvertTargetDaysToJSON converts TargetDays entity to JSON string
func ConvertTargetDaysToJSON(targetDays *entity.TargetDays) (*string, error) {
	if targetDays == nil {
		return nil, nil
	}

	targetDaysJSON, err := json.Marshal(targetDays)
	if err != nil {
		return nil, errors.New("failed to marshal target days")
	}

	targetDaysStr := string(targetDaysJSON)
	return &targetDaysStr, nil
}
