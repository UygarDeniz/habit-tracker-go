package entity

import (
	"errors"
	"slices"
	"strings"
	"time"
)

// TargetDays represents the days when a habit should be performed
type TargetDays struct {
	// For weekly habits: ["monday", "wednesday", "friday"]
	// For monthly habits: [1, 15] (days of month) or ["last"] for last day of month
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
			return err
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

var validFrequencies = []string{"daily", "weekly", "monthly"}

func isValidFrequency(frequency string) bool {
	return slices.Contains(validFrequencies, frequency)
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
			// Handle special case for "last" day of month
			if dayStr, ok := day.(string); ok {
				if dayStr == "last" {
					continue // "last" is valid for last day of any month
				}
				return errors.New("invalid month day: only numeric days (1-28) or 'last' are allowed")
			}

			// Handle numeric days
			dayNum, ok := day.(float64)
			dayInt := int(dayNum)
			if !ok || float64(dayInt) != dayNum {
				return errors.New("invalid month day: must be a number between 1-28 or 'last'")
			}

			// Restrict to days 1-28 to ensure they exist in all months
			if dayInt < 1 || dayInt > 28 {
				return errors.New("invalid month day: must be between 1-28 (to ensure availability in all months) or use 'last' for the last day of month")
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

// GetValidMonthlyDays returns the actual days for a given month/year, handling edge cases
func (td *TargetDays) GetValidMonthlyDays(year int, month time.Month) []int {
	if td == nil {
		return nil
	}

	var validDays []int
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for _, day := range td.Days {
		if dayStr, ok := day.(string); ok && dayStr == "last" {
			validDays = append(validDays, daysInMonth)
		} else if dayNum, ok := day.(float64); ok {
			dayInt := int(dayNum)
			if dayInt <= daysInMonth {
				validDays = append(validDays, dayInt)
			}
			// If the target day doesn't exist in this month, skip it
			// This prevents issues with habits set for day 31 in February
		}
	}

	return validDays
}

// IsValidForMonth checks if the target days are valid for a specific month
func (td *TargetDays) IsValidForMonth(year int, month time.Month) bool {
	if td == nil {
		return true
	}

	validDays := td.GetValidMonthlyDays(year, month)
	return len(validDays) > 0 // At least one day should be valid
}
