package services

import (
	"context"

	"github.com/deathcore1998/workout-tracker/internal/models"
	"github.com/deathcore1998/workout-tracker/internal/repositories"
)

type WorkoutService struct {
	repo *repositories.WorkoutRepository
}

func NewWorkoutService(repo *repositories.WorkoutRepository) *WorkoutService {
	return &WorkoutService{repo: repo}
}

func (s *WorkoutService) Create(ctx context.Context, workout *models.Workout) error {
	return s.repo.Create(ctx, workout)
}

func (s *WorkoutService) GetByUserID(ctx context.Context, userID int) ([]models.Workout, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *WorkoutService) GetByID(ctx context.Context, id, userID int) (*models.Workout, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *WorkoutService) Update(ctx context.Context, workout *models.Workout) error {
	return s.repo.Update(ctx, workout)
}

func (s *WorkoutService) Delete(ctx context.Context, id, userID int) error {
	return s.repo.Delete(ctx, id, userID)
}
