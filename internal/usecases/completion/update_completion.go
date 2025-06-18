package completion

import (
	"context"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/dto"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type UpdateCompletionUsecase struct {
	completionRepo repository.CompletionRepository
	habitRepo      repository.HabitRepository
}

func NewUpdateCompletionUsecase(completionRepo repository.CompletionRepository, habitRepo repository.HabitRepository) *UpdateCompletionUsecase {
	return &UpdateCompletionUsecase{
		completionRepo: completionRepo,
		habitRepo:      habitRepo,
	}
}

func (uc *UpdateCompletionUsecase) Execute(ctx context.Context, completionID, userID string, req dto.UpdateCompletionDTO) (*entity.HabitCompletion, error) {
	completion, err := uc.completionRepo.FindByID(ctx, completionID)
	if err != nil {
		return nil, err
	}

	if completion.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	originalCount := completion.Count

	if req.Count != nil {
		completion.Count = *req.Count
	}

	if req.Notes != nil {
		completion.SetNotes(*req.Notes)
	}

	if err := entity.ValidateCompletion(completion); err != nil {
		return nil, apperrors.ErrInvalidInput
	}

	habit, err := uc.habitRepo.FindByID(ctx, completion.HabitID)
	if err != nil {
		return nil, err
	}

	habit.TotalCompletions = habit.TotalCompletions - originalCount + completion.Count

	err = uc.completionRepo.Update(ctx, completion, habit)
	if err != nil {
		return nil, err
	}

	return completion, nil
}
