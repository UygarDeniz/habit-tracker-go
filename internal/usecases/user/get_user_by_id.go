package user

import (
	"context"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type GetUserByIDUsecase struct {
	userRepository repository.UserRepository
}

func NewGetUserByIDUsecase(userRepository repository.UserRepository) *GetUserByIDUsecase {
	return &GetUserByIDUsecase{userRepository: userRepository}
}

func (uc *GetUserByIDUsecase) Execute(ctx context.Context, userID string) (*entity.User, error) {
	user, err := uc.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}
