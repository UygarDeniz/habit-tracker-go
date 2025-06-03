package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	fmt.Println("Creating user in Postgres with context")
	return nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	fmt.Printf("Finding user by email in Postgres with context: %s\n", email)

	return nil, apperrors.ErrNotFound
}
