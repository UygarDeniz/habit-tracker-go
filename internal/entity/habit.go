package entity

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"time"
)

// TargetDays represents the days when a habit should be performed
type TargetDays struct {
	// For weekly habits: ["monday", "wednesday", "friday"]
	// For monthly habits: [1, 15] (days of month)
	Days []any `json:"days"`
}

type Habit struct {
	ID               string      `json:"id"`
	UserID           string      `json:"user_id"`
	Name             string      `json:"name"`
	Description      *string     `json:"description"`
	Motivation       *string     `json:"motivation"`
	Color            string      `json:"color"`
	Category         *string     `json:"category"`
	Frequency        string      `json:"frequency"`
	TargetCount      int         `json:"target_count"`
	TargetDays       *TargetDays `json:"target_days"`
	CurrentStreak    int         `json:"current_streak"`
	BestStreak       int         `json:"best_streak"`
	TotalCompletions int         `json:"total_completions"`
	IsActive         bool        `json:"is_active"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

func NewHabit(id, userID string, name, frequency string, targetCount int, description, motivation, category *string, targetDays *TargetDays, color string) (*Habit, error) {
	now := time.Now()
	habit := &Habit{
		ID:               id,
		UserID:           userID,
		Name:             name,
		Description:      description,
		Motivation:       motivation,
		Color:            color,
		Category:         category,
		Frequency:        frequency,
		TargetCount:      targetCount,
		TargetDays:       targetDays,
		CurrentStreak:    0,
		BestStreak:       0,
		TotalCompletions: 0,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := Validate(habit); err != nil {
		return nil, err
	}

	return habit, nil
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

func (h *Habit) SetTargetDays(days *TargetDays) error {
	if days != nil {
		if err := validateTargetDays(h.Frequency, days); err != nil {
			return errors.New("invalid target days")
		}
	}
	h.TargetDays = days
	return nil
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

func isValidColor(color string) bool {
	// Check if color is a valid hex color code
	return strings.HasPrefix(color, "#") && len(color) == 7
}

func validateTargetDays(frequency string, days *TargetDays) error {
	if days == nil {
		return nil
	}
	validDays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

	switch frequency {
	case "daily":
		return nil
	case "weekly":

		for _, day := range days.Days {
			dayStr, ok := day.(string)
			if !ok || !slices.Contains(validDays, dayStr) {
				return errors.New("invalid weekday")
			}
		}
	case "monthly":
		for _, day := range days.Days {
			dayNum, ok := day.(float64)
			if !ok || dayNum < 1 || dayNum > 31 {
				return errors.New("invalid month day")
			}
		}
	default:
		return errors.New("invalid frequency")
	}

	return nil
}

func Validate(habit *Habit) error {
	if habit.ID == "" {
		return errors.New("id is required")
	}
	if habit.UserID == "" {
		return errors.New("user ID is required")
	}
	if strings.TrimSpace(habit.Name) == "" {
		return errors.New("name is required")
	}
	if !isValidFrequency(habit.Frequency) {
		return errors.New("invalid frequency")
	}
	if habit.TargetCount <= 0 {
		return errors.New("target count must be positive")
	}
	if !isValidColor(habit.Color) {
		return errors.New("invalid color format")
	}
	if err := validateTargetDays(habit.Frequency, habit.TargetDays); err != nil {
		return err
	}

	return nil
}

// ConvertTargetDaysFromJSON converts JSON string to TargetDays entity
func ConvertTargetDaysFromJSON(targetDaysJSON *string) (*TargetDays, error) {
	if targetDaysJSON == nil || *targetDaysJSON == "" {
		return nil, nil
	}

	var targetDaysData struct {
		Days []any `json:"days"`
	}

	if err := json.Unmarshal([]byte(*targetDaysJSON), &targetDaysData); err != nil {
		return nil, errors.New("invalid target days JSON format")
	}

	return &TargetDays{Days: targetDaysData.Days}, nil
}

// ConvertTargetDaysToJSON converts TargetDays entity to JSON string
func ConvertTargetDaysToJSON(targetDays *TargetDays) (*string, error) {
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
