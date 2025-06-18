package habit

import (
	"context"
	"time"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/dto"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type UpdateHabitUsecase struct {
	habitRepository repository.HabitRepository
}

func NewUpdateHabitUsecase(habitRepository repository.HabitRepository) *UpdateHabitUsecase {
	return &UpdateHabitUsecase{habitRepository: habitRepository}
}

func (uc *UpdateHabitUsecase) Execute(ctx context.Context, habitID string, userID string, req dto.UpdateHabitDTO) (*entity.Habit, error) {
	habit, err := uc.habitRepository.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	if habit.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	// Apply updates from DTO to entity
	if req.Name != nil {
		habit.Name = *req.Name
	}
	if req.Description != nil {
		habit.SetDescription(*req.Description)
	}
	if req.Motivation != nil {
		habit.SetMotivation(*req.Motivation)
	}
	if req.Color != nil {
		habit.Color = *req.Color
	}
	if req.Category != nil {
		habit.SetCategory(*req.Category)
	}
	if req.Frequency != nil {
		habit.Frequency = *req.Frequency
	}
	if req.TargetCount != nil {
		habit.TargetCount = *req.TargetCount
	}

	if req.TargetDays != nil {
		targetDays, err := dto.ConvertTargetDaysFromJSON(req.TargetDays)
		if err != nil {
			return nil, apperrors.ErrInvalidInput
		}
		if err := habit.SetTargetDays(targetDays); err != nil {
			return nil, apperrors.ErrInvalidInput
		}
	} else if req.Frequency != nil {
		// Re-validate target days if frequency is changed but target days are not
		if err := habit.SetTargetDays(habit.TargetDays); err != nil {
			return nil, apperrors.ErrInvalidInput
		}
	}

	if req.IsActive != nil {
		if *req.IsActive {
			habit.Activate()
		} else {
			habit.Deactivate()
		}
	}

	habit.UpdatedAt = time.Now()

	// Re-validate the entity after updates
	if err := entity.Validate(habit); err != nil {
		return nil, apperrors.ErrInvalidInput
	}

	err = uc.habitRepository.Update(ctx, habit)
	if err != nil {
		return nil, err
	}

	return habit, nil
}
