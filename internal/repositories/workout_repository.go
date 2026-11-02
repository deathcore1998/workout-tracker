package repositories

import (
	"context"
	"fmt"

	"github.com/deathcore1998/workout-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutRepository struct {
	db *pgxpool.Pool
}

func NewWorkoutRepository(db *pgxpool.Pool) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

func (repo *WorkoutRepository) Create(ctx context.Context, workout *models.Workout) error {
	query := `INSERT INTO workouts (user_id, name, date, notes) VALUES($1,$2,$3,$4) 
	RETURNING id, created_at, updated_at`
	row := repo.db.QueryRow(ctx, query, workout.UserID, workout.Name, workout.Date, workout.Notes)
	return row.Scan(&workout.ID, &workout.CreatedAt, &workout.UpdatedAt)
}

func (repo *WorkoutRepository) GetByUserID(ctx context.Context, userID int) ([]models.Workout, error) {
	query := `SELECT id, user_id, name, date, notes, created_at, updated_at
			FROM workouts WHERE user_id = $1 ORDER BY date DESC`
	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var workouts []models.Workout
	for rows.Next() {
		var workout models.Workout
		err := rows.Scan(&workout.ID, &workout.UserID, &workout.Name, &workout.Date, &workout.Notes, &workout.CreatedAt, &workout.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return workouts, nil
}

func (repo *WorkoutRepository) GetByID(ctx context.Context, id, userID int) (*models.Workout, error) {
	query := `SELECT id, user_id, name, date, notes, created_at, updated_at
	FROM workouts WHERE id = $1 AND user_id = $2`
	row := repo.db.QueryRow(ctx, query, id, userID)
	workout := &models.Workout{}
	err := row.Scan(&workout.ID, &workout.UserID, &workout.Name, &workout.Date, &workout.Notes, &workout.CreatedAt, &workout.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (repo *WorkoutRepository) Update(ctx context.Context, workout *models.Workout) error {
	query := `UPDATE workouts
	SET name = $1, date = $2, notes = $3, updated_at = NOW()
	WHERE id = $4 AND user_id = $5
	RETURNING updated_at`
	row := repo.db.QueryRow(ctx, query, workout.Name, workout.Date, workout.Notes, workout.ID, workout.UserID)
	return row.Scan(&workout.UpdatedAt)
}

func (repo *WorkoutRepository) Delete(ctx context.Context, id, userID int) error {
	query := `DELETE FROM workouts WHERE id = $1 AND user_id = $2`
	result, err := repo.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("workout not found")
	}
	return nil
}
