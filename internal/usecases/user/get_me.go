package user

import (
	"context"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/middleware"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type GetMeUsecase struct {
	userRepository repository.UserRepository
}

func NewGetMeUsecase(userRepository repository.UserRepository) *GetMeUsecase {
	return &GetMeUsecase{userRepository: userRepository}
}

func (uc *GetMeUsecase) Execute(ctx context.Context) (*entity.User, error) {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	user, err := uc.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}
