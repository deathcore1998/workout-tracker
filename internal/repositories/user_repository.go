package repositories

import (
	"context"

	"github.com/deathcore1998/workout-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	queryRow := r.db.QueryRow(ctx, query, user.Email, user.Password, user.Name)
	return queryRow.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}
