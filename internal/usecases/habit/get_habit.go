package habit

import (
	"context"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type GetHabitUsecase struct {
	habitRepository repository.HabitRepository
}

func NewGetHabitUsecase(habitRepository repository.HabitRepository) *GetHabitUsecase {
	return &GetHabitUsecase{habitRepository: habitRepository}
}

func (uc *GetHabitUsecase) Execute(ctx context.Context, habitID string, userID string) (*entity.Habit, error) {
	habit, err := uc.habitRepository.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	if habit.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	return habit, nil
}
