package completion

import (
	"context"
	"time"

	"github.com/uygardeniz/habit-tracker/internal/dto"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type GetCompletionsUsecase struct {
	completionRepo repository.CompletionRepository
}

func NewGetCompletionsUsecase(completionRepo repository.CompletionRepository) *GetCompletionsUsecase {
	return &GetCompletionsUsecase{
		completionRepo: completionRepo,
	}
}

func (uc *GetCompletionsUsecase) Execute(ctx context.Context, userID string, query dto.GetCompletionsQueryDTO) ([]*entity.HabitCompletion, error) {
	var startDate, endDate *time.Time
	var err error

	if query.StartDate != nil && *query.StartDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", *query.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = &parsedStartDate
	}

	if query.EndDate != nil && *query.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", *query.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &parsedEndDate
	}

	limit := 50
	if query.Limit != nil {
		limit = *query.Limit
	}

	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
	}

	completions, err := uc.completionRepo.FindByUserID(ctx, userID, query.HabitID, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}

	return completions, nil
}
