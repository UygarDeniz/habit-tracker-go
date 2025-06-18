package completion

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/dto"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type CreateCompletionUsecase struct {
	completionRepo repository.CompletionRepository
	habitRepo      repository.HabitRepository
}

func NewCreateCompletionUsecase(completionRepo repository.CompletionRepository, habitRepo repository.HabitRepository) *CreateCompletionUsecase {
	return &CreateCompletionUsecase{
		completionRepo: completionRepo,
		habitRepo:      habitRepo,
	}
}

func (uc *CreateCompletionUsecase) Execute(ctx context.Context, habitID, userID string, req dto.CreateCompletionDTO) (*entity.HabitCompletion, error) {
	habit, err := uc.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	if habit.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	completionDate, err := time.Parse("2006-01-02", req.CompletionDate)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}

	existingCompletion, err := uc.completionRepo.FindByHabitIDAndDate(ctx, habitID, completionDate)
	if err != nil && err != apperrors.ErrNotFound {
		return nil, err
	}

	if existingCompletion != nil {
		return nil, apperrors.ErrAlreadyExists
	}

	completionID := uuid.New().String()
	completion, err := entity.NewHabitCompletion(completionID, habitID, userID, completionDate, req.Count, req.Notes)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}

	habit.IncrementCompletions()
	if shouldIncrementStreak(completionDate) {
		habit.IncrementStreak()
	}

	return uc.completionRepo.Create(ctx, completion, habit)
}

func shouldIncrementStreak(completionDate time.Time) bool {
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	completionDay := completionDate.Truncate(24 * time.Hour)

	return completionDay.Equal(today) || completionDay.Equal(yesterday)
}
