package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/entity"
)

type CompletionRepository interface {
	Create(ctx context.Context, completion *entity.HabitCompletion, habit *entity.Habit) (*entity.HabitCompletion, error)
	FindByID(ctx context.Context, id string) (*entity.HabitCompletion, error)
	FindByUserID(ctx context.Context, userID string, habitID *string, startDate, endDate *time.Time, limit, offset int) ([]*entity.HabitCompletion, error)
	FindByHabitID(ctx context.Context, habitID string, startDate, endDate *time.Time, limit, offset int) ([]*entity.HabitCompletion, error)
	FindByHabitIDAndDate(ctx context.Context, habitID string, date time.Time) (*entity.HabitCompletion, error)
	Update(ctx context.Context, completion *entity.HabitCompletion, habit *entity.Habit) error
	Delete(ctx context.Context, id string, habit *entity.Habit) error
	CountByUserID(ctx context.Context, userID string, habitID *string, startDate, endDate *time.Time) (int, error)
}

type PostgresCompletionRepository struct {
	db *sql.DB
}

func NewPostgresCompletionRepository(db *sql.DB) CompletionRepository {
	return &PostgresCompletionRepository{db: db}
}

func (r *PostgresCompletionRepository) Create(ctx context.Context, completion *entity.HabitCompletion, habit *entity.Habit) (*entity.HabitCompletion, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	completionQuery := `
		INSERT INTO habit_completions (id, habit_id, user_id, completed_at, completion_date, count, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, habit_id, user_id, completed_at, completion_date, count, notes, created_at
	`

	row := tx.QueryRowContext(ctx, completionQuery,
		completion.ID, completion.HabitID, completion.UserID, completion.CompletedAt,
		completion.CompletionDate, completion.Count, completion.Notes, completion.CreatedAt,
	)

	var createdCompletion entity.HabitCompletion
	err = row.Scan(
		&createdCompletion.ID, &createdCompletion.HabitID, &createdCompletion.UserID,
		&createdCompletion.CompletedAt, &createdCompletion.CompletionDate,
		&createdCompletion.Count, &createdCompletion.Notes, &createdCompletion.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	habitQuery := `
		UPDATE habits
		SET current_streak = $1, best_streak = $2, total_completions = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := tx.ExecContext(ctx, habitQuery, habit.CurrentStreak, habit.BestStreak, habit.TotalCompletions, habit.UpdatedAt, habit.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &createdCompletion, nil
}

func (r *PostgresCompletionRepository) Delete(ctx context.Context, id string, habit *entity.Habit) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM habit_completions WHERE id = $1`
	result, err := tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	habitQuery := `
		UPDATE habits
		SET total_completions = $1
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, habitQuery, habit.TotalCompletions, habit.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresCompletionRepository) FindByID(ctx context.Context, id string) (*entity.HabitCompletion, error) {
	query := `
		SELECT id, habit_id, user_id, completed_at, completion_date, count, notes, created_at
		FROM habit_completions
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var completion entity.HabitCompletion
	err := row.Scan(
		&completion.ID, &completion.HabitID, &completion.UserID,
		&completion.CompletedAt, &completion.CompletionDate,
		&completion.Count, &completion.Notes, &completion.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return &completion, nil
}

func (r *PostgresCompletionRepository) FindByUserID(ctx context.Context, userID string, habitID *string, startDate, endDate *time.Time, limit, offset int) ([]*entity.HabitCompletion, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if habitID != nil {
		conditions = append(conditions, fmt.Sprintf("habit_id = $%d", argIndex))
		args = append(args, *habitID)
		argIndex++
	}

	if startDate != nil {
		conditions = append(conditions, fmt.Sprintf("completion_date >= $%d", argIndex))
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		conditions = append(conditions, fmt.Sprintf("completion_date <= $%d", argIndex))
		args = append(args, *endDate)
		argIndex++
	}

	query := fmt.Sprintf(`
		SELECT id, habit_id, user_id, completed_at, completion_date, count, notes, created_at
		FROM habit_completions
		WHERE %s
		ORDER BY completion_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var completions []*entity.HabitCompletion
	for rows.Next() {
		var completion entity.HabitCompletion
		err := rows.Scan(
			&completion.ID, &completion.HabitID, &completion.UserID,
			&completion.CompletedAt, &completion.CompletionDate,
			&completion.Count, &completion.Notes, &completion.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		completions = append(completions, &completion)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return completions, nil
}

func (r *PostgresCompletionRepository) FindByHabitID(ctx context.Context, habitID string, startDate, endDate *time.Time, limit, offset int) ([]*entity.HabitCompletion, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("habit_id = $%d", argIndex))
	args = append(args, habitID)
	argIndex++

	if startDate != nil {
		conditions = append(conditions, fmt.Sprintf("completion_date >= $%d", argIndex))
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		conditions = append(conditions, fmt.Sprintf("completion_date <= $%d", argIndex))
		args = append(args, *endDate)
		argIndex++
	}

	query := fmt.Sprintf(`
		SELECT id, habit_id, user_id, completed_at, completion_date, count, notes, created_at
		FROM habit_completions
		WHERE %s
		ORDER BY completion_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var completions []*entity.HabitCompletion
	for rows.Next() {
		var completion entity.HabitCompletion
		err := rows.Scan(
			&completion.ID, &completion.HabitID, &completion.UserID,
			&completion.CompletedAt, &completion.CompletionDate,
			&completion.Count, &completion.Notes, &completion.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		completions = append(completions, &completion)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return completions, nil
}

func (r *PostgresCompletionRepository) FindByHabitIDAndDate(ctx context.Context, habitID string, date time.Time) (*entity.HabitCompletion, error) {
	query := `
		SELECT id, habit_id, user_id, completed_at, completion_date, count, notes, created_at
		FROM habit_completions
		WHERE habit_id = $1 AND completion_date = $2
	`
	row := r.db.QueryRowContext(ctx, query, habitID, date)

	var completion entity.HabitCompletion
	err := row.Scan(
		&completion.ID, &completion.HabitID, &completion.UserID,
		&completion.CompletedAt, &completion.CompletionDate,
		&completion.Count, &completion.Notes, &completion.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return &completion, nil
}

func (r *PostgresCompletionRepository) Update(ctx context.Context, completion *entity.HabitCompletion, habit *entity.Habit) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update completion
	completionQuery := `
		UPDATE habit_completions
		SET count = $1, notes = $2
		WHERE id = $3
	`
	_, err = tx.ExecContext(ctx, completionQuery, completion.Count, completion.Notes, completion.ID)
	if err != nil {
		return err
	}

	// Update habit statistics
	habitQuery := `
		UPDATE habits
		SET total_completions = $1
		WHERE id = $2
	`
	result, err := tx.ExecContext(ctx, habitQuery, habit.TotalCompletions, habit.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return tx.Commit()
}

func (r *PostgresCompletionRepository) CountByUserID(ctx context.Context, userID string, habitID *string, startDate, endDate *time.Time) (int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if habitID != nil {
		conditions = append(conditions, fmt.Sprintf("habit_id = $%d", argIndex))
		args = append(args, *habitID)
		argIndex++
	}

	if startDate != nil {
		conditions = append(conditions, fmt.Sprintf("completion_date >= $%d", argIndex))
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		conditions = append(conditions, fmt.Sprintf("completion_date <= $%d", argIndex))
		args = append(args, *endDate)
		argIndex++
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM habit_completions
		WHERE %s
	`, strings.Join(conditions, " AND "))

	row := r.db.QueryRowContext(ctx, query, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
