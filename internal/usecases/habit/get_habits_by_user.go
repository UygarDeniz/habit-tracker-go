package habit

import (
	"context"

	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type GetHabitsByUserUsecase struct {
	habitRepository repository.HabitRepository
}

func NewGetHabitsByUserUsecase(habitRepository repository.HabitRepository) *GetHabitsByUserUsecase {
	return &GetHabitsByUserUsecase{habitRepository: habitRepository}
}

func (uc *GetHabitsByUserUsecase) Execute(ctx context.Context, userID string) ([]*entity.Habit, error) {
	return uc.habitRepository.FindByUserID(ctx, userID)
}
