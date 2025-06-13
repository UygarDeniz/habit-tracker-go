package habit

import (
	"context"

	"github.com/google/uuid"
	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/dto"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type CreateHabitUsecase struct {
	habitRepository repository.HabitRepository
}

func NewCreateHabitUsecase(habitRepository repository.HabitRepository) *CreateHabitUsecase {
	return &CreateHabitUsecase{habitRepository: habitRepository}
}

func (uc *CreateHabitUsecase) Execute(ctx context.Context, userID string, req dto.CreateHabitDTO) (*entity.Habit, error) {
	targetDays, err := entity.ConvertTargetDaysFromJSON(req.TargetDays)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}

	habit, err := entity.NewHabit(uuid.New().String(), userID, req.Name, req.Frequency,
		req.TargetCount, req.Description, req.Motivation, req.Category, targetDays, req.Color)

	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}

	habit, err = uc.habitRepository.Create(ctx, habit)
	if err != nil {
		return nil, err
	}
	return habit, nil
}
