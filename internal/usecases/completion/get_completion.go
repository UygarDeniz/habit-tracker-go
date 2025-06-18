package completion

import (
	"context"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type GetCompletionUsecase struct {
	completionRepo repository.CompletionRepository
}

func NewGetCompletionUsecase(completionRepo repository.CompletionRepository) *GetCompletionUsecase {
	return &GetCompletionUsecase{
		completionRepo: completionRepo,
	}
}

func (uc *GetCompletionUsecase) Execute(ctx context.Context, completionID, userID string) (*entity.HabitCompletion, error) {
	completion, err := uc.completionRepo.FindByID(ctx, completionID)
	if err != nil {
		return nil, err
	}

	if completion.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	return completion, nil
}
