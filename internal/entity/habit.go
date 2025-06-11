package entity

import (
	"errors"
	"strings"
	"time"
)

type Habit struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	Description      *string   `json:"description"`
	Motivation       *string   `json:"motivation"`
	Color            string    `json:"color"`
	Category         *string   `json:"category"`
	Frequency        string    `json:"frequency"`
	TargetCount      int       `json:"target_count"`
	TargetDays       *string   `json:"target_days"`
	CurrentStreak    int       `json:"current_streak"`
	BestStreak       int       `json:"best_streak"`
	TotalCompletions int       `json:"total_completions"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewHabit(userID string, name, frequency string, targetCount int) (*Habit, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("habit name cannot be empty")
	}
	if userID == "" {
		return nil, errors.New("user ID must be positive")
	}
	if targetCount <= 0 {
		return nil, errors.New("target count must be positive")
	}
	if !isValidFrequency(frequency) {
		return nil, errors.New("frequency must be one of: daily, weekly, monthly, custom")
	}

	now := time.Now()
	return &Habit{
		UserID:           userID,
		Name:             name,
		Color:            "#3B82F6",
		Frequency:        frequency,
		TargetCount:      targetCount,
		CurrentStreak:    0,
		BestStreak:       0,
		TotalCompletions: 0,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

func (h *Habit) SetDescription(description string) {
	if strings.TrimSpace(description) == "" {
		h.Description = nil
	} else {
		h.Description = &description
	}
}

func (h *Habit) SetMotivation(motivation string) {
	if strings.TrimSpace(motivation) == "" {
		h.Motivation = nil
	} else {
		h.Motivation = &motivation
	}
}

func (h *Habit) SetCategory(category string) {
	if strings.TrimSpace(category) == "" {
		h.Category = nil
	} else {
		h.Category = &category
	}
}

func (h *Habit) IncrementStreak() {
	h.CurrentStreak++
	if h.CurrentStreak > h.BestStreak {
		h.BestStreak = h.CurrentStreak
	}
}

func (h *Habit) ResetStreak() {
	h.CurrentStreak = 0
}

func (h *Habit) IncrementCompletions() {
	h.TotalCompletions++
}

func (h *Habit) Deactivate() {
	h.IsActive = false
}

func (h *Habit) Activate() {
	h.IsActive = true
}

func isValidFrequency(frequency string) bool {
	validFrequencies := []string{"daily", "weekly", "monthly", "custom"}
	for _, valid := range validFrequencies {
		if frequency == valid {
			return true
		}
	}
	return false
}
