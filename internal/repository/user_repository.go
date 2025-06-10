package repository

import (
	"context"
	"database/sql"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, name, picture, google_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name, user.Picture, user.GoogleID)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	query := `
		SELECT id, email, name, picture, google_id
		FROM users
		WHERE google_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, googleID)

	var foundUser entity.User

	err := row.Scan(
		&foundUser.ID,
		&foundUser.Email,
		&foundUser.Name,
		&foundUser.Picture,
		&foundUser.GoogleID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return &foundUser, nil
}
