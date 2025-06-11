package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
)

type LoginOrRegisterGoogleUserUsecase struct {
	userRepo repository.UserRepository
}

func NewLoginOrRegisterGoogleUserUsecase(userRepo repository.UserRepository) *LoginOrRegisterGoogleUserUsecase {
	return &LoginOrRegisterGoogleUserUsecase{userRepo: userRepo}
}

func (uc *LoginOrRegisterGoogleUserUsecase) Execute(ctx context.Context, googleID, email, name, picture string) (*entity.User, error) {
	user, err := uc.userRepo.FindByGoogleID(ctx, googleID)
	if err != nil {
		if err == apperrors.ErrNotFound {
			// User not found, create a new one
			newUser, err := entity.NewUser(uuid.NewString(), email, name, picture, googleID)
			if err != nil {
				return nil, err
			}
			err = uc.userRepo.Create(ctx, newUser)
			if err != nil {
				return nil, err
			}
			return newUser, nil
		}
		return nil, err
	}

	return user, nil
}
