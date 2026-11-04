package services

import (
	"context"

	"github.com/deathcore1998/workout-tracker/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
}

type WorkoutRepository interface {
	Create(ctx context.Context, workout *models.Workout) error
	GetByUserID(ctx context.Context, userID int) ([]models.Workout, error)
	GetByID(ctx context.Context, id, userID int) (*models.Workout, error)
	Update(ctx context.Context, workout *models.Workout) error
	Delete(ctx context.Context, id, userID int) error
}
