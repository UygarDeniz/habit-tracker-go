package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
	"github.com/uygardeniz/habit-tracker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserUsecase struct {
	userRepository repository.UserRepository
}

func NewCreateUserUsecase(userRepository repository.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{userRepository: userRepository}
}

func (u *CreateUserUsecase) Execute(ctx context.Context, email, username, password string) (*entity.User, error) {

	existingUser, err := u.userRepository.FindByEmail(ctx, email)

	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, fmt.Errorf("failed to check for existing user: %w", err)
	}

	if existingUser != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := entity.NewUser(uuid.New().String(), email, username, string(hashedPassword))

	if err != nil {
		return nil, fmt.Errorf("failed to create new user entity: %w", err)
	}

	if err = u.userRepository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
