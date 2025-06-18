package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
)

type HabitRepository interface {
	Create(ctx context.Context, habit *entity.Habit) (*entity.Habit, error)
	FindByID(ctx context.Context, id string) (*entity.Habit, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.Habit, error)
	Update(ctx context.Context, habit *entity.Habit) error
	Delete(ctx context.Context, id string) error
}

type PostgresHabitRepository struct {
	db *sql.DB
}

func NewPostgresHabitRepository(db *sql.DB) HabitRepository {
	return &PostgresHabitRepository{db: db}
}

func (r *PostgresHabitRepository) Create(ctx context.Context, habit *entity.Habit) (*entity.Habit, error) {
	query := `
		INSERT INTO habits (id, user_id, name, description, motivation, color, category, frequency, target_count, target_days, current_streak, best_streak, total_completions, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, user_id, name, description, motivation, color, category, frequency, target_count, target_days, current_streak, best_streak, total_completions, is_active, created_at, updated_at
	`

	var targetDaysJSON []byte
	if habit.TargetDays != nil {
		var err error
		targetDaysJSON, err = json.Marshal(habit.TargetDays)
		if err != nil {
			return nil, err
		}
	}

	row := r.db.QueryRowContext(ctx, query,
		habit.ID, habit.UserID, habit.Name, habit.Description, habit.Motivation,
		habit.Color, habit.Category, habit.Frequency, habit.TargetCount,
		targetDaysJSON, habit.CurrentStreak, habit.BestStreak,
		habit.TotalCompletions, habit.IsActive, habit.CreatedAt, habit.UpdatedAt,
	)

	var createdHabit entity.Habit
	var targetDaysBytes []byte

	err := row.Scan(
		&createdHabit.ID, &createdHabit.UserID, &createdHabit.Name,
		&createdHabit.Description, &createdHabit.Motivation, &createdHabit.Color,
		&createdHabit.Category, &createdHabit.Frequency, &createdHabit.TargetCount,
		&targetDaysBytes, &createdHabit.CurrentStreak, &createdHabit.BestStreak,
		&createdHabit.TotalCompletions, &createdHabit.IsActive,
		&createdHabit.CreatedAt, &createdHabit.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Handle target_days JSON conversion
	if targetDaysBytes != nil {
		var targetDays entity.TargetDays
		if err := json.Unmarshal(targetDaysBytes, &targetDays); err != nil {
			return nil, err
		}
		createdHabit.TargetDays = &targetDays
	}

	return &createdHabit, nil
}

func (r *PostgresHabitRepository) FindByID(ctx context.Context, id string) (*entity.Habit, error) {
	query := `
		SELECT id, user_id, name, description, motivation, color, category, frequency, target_count, target_days, current_streak, best_streak, total_completions, is_active, created_at, updated_at
		FROM habits
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var habit entity.Habit
	var targetDaysBytes []byte

	err := row.Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Description,
		&habit.Motivation,
		&habit.Color,
		&habit.Category,
		&habit.Frequency,
		&habit.TargetCount,
		&targetDaysBytes,
		&habit.CurrentStreak,
		&habit.BestStreak,
		&habit.TotalCompletions,
		&habit.IsActive,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	if targetDaysBytes != nil {
		var targetDays entity.TargetDays
		if err := json.Unmarshal(targetDaysBytes, &targetDays); err != nil {
			return nil, err
		}
		habit.TargetDays = &targetDays
	}

	return &habit, nil
}

func (r *PostgresHabitRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Habit, error) {
	query := `
		SELECT id, user_id, name, description, motivation, color, category, frequency, target_count, target_days, current_streak, best_streak, total_completions, is_active, created_at, updated_at
		FROM habits
		WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []*entity.Habit
	for rows.Next() {
		var habit entity.Habit
		var targetDaysBytes []byte
		err := rows.Scan(
			&habit.ID,
			&habit.UserID,
			&habit.Name,
			&habit.Description,
			&habit.Motivation,
			&habit.Color,
			&habit.Category,
			&habit.Frequency,
			&habit.TargetCount,
			&targetDaysBytes,
			&habit.CurrentStreak,
			&habit.BestStreak,
			&habit.TotalCompletions,
			&habit.IsActive,
			&habit.CreatedAt,
			&habit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if targetDaysBytes != nil {
			var targetDays entity.TargetDays
			if err := json.Unmarshal(targetDaysBytes, &targetDays); err != nil {
				return nil, err
			}
			habit.TargetDays = &targetDays
		}

		habits = append(habits, &habit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return habits, nil
}

func (r *PostgresHabitRepository) Update(ctx context.Context, habit *entity.Habit) error {
	query := `
		UPDATE habits
		SET name = $1, description = $2, motivation = $3, color = $4, category = $5, frequency = $6, target_count = $7, target_days = $8, current_streak = $9, best_streak = $10, total_completions = $11, is_active = $12, updated_at = $13
		WHERE id = $14
	`

	var targetDaysJSON []byte
	var err error
	if habit.TargetDays != nil {
		targetDaysJSON, err = json.Marshal(habit.TargetDays)
		if err != nil {
			return err
		}
	}

	result, err := r.db.ExecContext(ctx, query, habit.Name, habit.Description, habit.Motivation, habit.Color, habit.Category, habit.Frequency, habit.TargetCount, targetDaysJSON, habit.CurrentStreak, habit.BestStreak, habit.TotalCompletions, habit.IsActive, habit.UpdatedAt, habit.ID)

	if err != nil {
		return err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsEffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *PostgresHabitRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM habits
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}
