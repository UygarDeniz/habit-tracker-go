package habit

import (
	"context"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type DeleteHabitUsecase struct {
	habitRepository repository.HabitRepository
}

func NewDeleteHabitUsecase(habitRepository repository.HabitRepository) *DeleteHabitUsecase {
	return &DeleteHabitUsecase{habitRepository: habitRepository}
}

func (uc *DeleteHabitUsecase) Execute(ctx context.Context, habitID string, userID string) error {
	// First check if habit exists and user owns it
	habit, err := uc.habitRepository.FindByID(ctx, habitID)
	if err != nil {
		return err
	}

	if habit.UserID != userID {
		return apperrors.ErrForbidden
	}

	// Delete the habit
	err = uc.habitRepository.Delete(ctx, habitID)
	if err != nil {
		return err
	}

	return nil
}
