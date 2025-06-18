package entity

import (
	"errors"
	"strings"
	"time"
)

// HabitCompletion represents a completion record for a habit
type HabitCompletion struct {
	ID             string    `json:"id"`
	HabitID        string    `json:"habit_id"`
	UserID         string    `json:"user_id"`
	CompletedAt    time.Time `json:"completed_at"`
	CompletionDate time.Time `json:"completion_date"`
	Count          int       `json:"count"`
	Notes          *string   `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
}

// NewHabitCompletion creates a new habit completion record
func NewHabitCompletion(id, habitID, userID string, completionDate time.Time, count int, notes *string) (*HabitCompletion, error) {
	now := time.Now()

	completion := &HabitCompletion{
		ID:             id,
		HabitID:        habitID,
		UserID:         userID,
		CompletedAt:    now,
		CompletionDate: completionDate,
		Count:          count,
		Notes:          notes,
		CreatedAt:      now,
	}

	if err := ValidateCompletion(completion); err != nil {
		return nil, err
	}

	return completion, nil
}

// SetNotes sets the notes for the completion
func (c *HabitCompletion) SetNotes(notes string) {
	if strings.TrimSpace(notes) == "" {
		c.Notes = nil
	} else {
		c.Notes = &notes
	}
}

// ValidateCompletion validates a habit completion
func ValidateCompletion(completion *HabitCompletion) error {
	if completion.ID == "" {
		return errors.New("id is required")
	}
	if completion.HabitID == "" {
		return errors.New("habit ID is required")
	}
	if completion.UserID == "" {
		return errors.New("user ID is required")
	}
	if completion.Count <= 0 {
		return errors.New("count must be positive")
	}
	if completion.CompletionDate.IsZero() {
		return errors.New("completion date is required")
	}

	return nil
}
