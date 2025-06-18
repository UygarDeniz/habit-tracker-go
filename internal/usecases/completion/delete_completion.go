package completion

import (
	"context"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type DeleteCompletionUsecase struct {
	completionRepo repository.CompletionRepository
	habitRepo      repository.HabitRepository
}

func NewDeleteCompletionUsecase(completionRepo repository.CompletionRepository, habitRepo repository.HabitRepository) *DeleteCompletionUsecase {
	return &DeleteCompletionUsecase{
		completionRepo: completionRepo,
		habitRepo:      habitRepo,
	}
}

func (uc *DeleteCompletionUsecase) Execute(ctx context.Context, completionID, userID string) error {
	completion, err := uc.completionRepo.FindByID(ctx, completionID)
	if err != nil {
		return err
	}

	if completion.UserID != userID {
		return apperrors.ErrForbidden
	}

	habit, err := uc.habitRepo.FindByID(ctx, completion.HabitID)
	if err != nil {
		return err
	}

	if habit.TotalCompletions > 0 {
		habit.TotalCompletions--
	}

	return uc.completionRepo.Delete(ctx, completionID, habit)
}
