package services

import (
	"context"
	"testing"
	"time"

	"github.com/deathcore1998/workout-tracker/internal/models"
)

type mockWorkoutRepository struct {
	workouts map[int]*models.Workout
	nextID   int
}

func newMockWorkoutRepository() *mockWorkoutRepository {
	return &mockWorkoutRepository{workouts: make(map[int]*models.Workout), nextID: 1}
}

func (mock *mockWorkoutRepository) Create(ctx context.Context, workout *models.Workout) error {
	workout.ID = mock.nextID
	mock.nextID++
	workout.CreatedAt = time.Now()
	workout.UpdatedAt = time.Now()
	mock.workouts[workout.ID] = workout
	return nil
}

func (mock *mockWorkoutRepository) GetByUserID(ctx context.Context, userID int) ([]models.Workout, error) {
	var workouts []models.Workout
	for _, workout := range mock.workouts {
		if workout.UserID == userID {
			workouts = append(workouts, *workout)
		}
	}
	return workouts, nil
}

func (mock *mockWorkoutRepository) GetByID(ctx context.Context, id, userID int) (*models.Workout, error) {
	existing, ok := mock.workouts[id]
	if ok && existing.UserID == userID {
		return existing, nil
	}
	return nil, nil
}

func (mock *mockWorkoutRepository) Update(ctx context.Context, workout *models.Workout) error {
	existing, ok := mock.workouts[workout.ID]
	if ok && existing.UserID == workout.UserID {
		workout.CreatedAt = existing.CreatedAt
		workout.UpdatedAt = time.Now()
		mock.workouts[workout.ID] = workout
	}
	return nil
}

func (mock *mockWorkoutRepository) Delete(ctx context.Context, id, userID int) error {
	existing, ok := mock.workouts[id]
	if ok && existing.UserID == userID {
		delete(mock.workouts, id)
	}
	return nil
}

func createTestWorkout(t *testing.T) *models.Workout {
	t.Helper()

	return &models.Workout{
		UserID: 1,
		Name:   "Chest",
		Date:   time.Now(),
		Notes:  "Hard",
	}
}

func TestWorkoutService_Create(t *testing.T) {
	repo := newMockWorkoutRepository()
	service := NewWorkoutService(repo)

	testWorkout := createTestWorkout(t)
	err := service.Create(context.Background(), testWorkout)
	if err != nil {
		t.Errorf("Create failed %v", err)
	}

	if testWorkout.ID == 0 {
		t.Errorf("expected workout ID, got 0")
	}

	if testWorkout.UserID != 1 {
		t.Errorf("expected userID 1, got %d", testWorkout.UserID)
	}

	if testWorkout.Name != "Chest" {
		t.Errorf("expected name Chest, got %s", testWorkout.Name)
	}
}

func TestWorkoutService_GetByUserID(t *testing.T) {
	repo := newMockWorkoutRepository()
	service := NewWorkoutService(repo)

	testWorkout := createTestWorkout(t)

	err := service.Create(context.Background(), testWorkout)
	if err != nil {
		t.Errorf("Create failed %v", err)
	}

	workouts, err := service.GetByUserID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}

	if len(workouts) != 1 {
		t.Fatalf("expected 1 workout, got %d", len(workouts))
	}

	if workouts[0].Name != "Chest" {
		t.Errorf("expected name Chest, got %s", workouts[0].Name)
	}
}

func TestWorkoutService_GetByID(t *testing.T) {
	repo := newMockWorkoutRepository()
	service := NewWorkoutService(repo)

	testWorkout := createTestWorkout(t)

	err := service.Create(context.Background(), testWorkout)
	if err != nil {
		t.Errorf("Create failed %v", err)
	}

	workout, err := service.GetByID(context.Background(), testWorkout.ID, 1)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if workout == nil {
		t.Fatalf("workout is nil")
	}

	if workout.Name != testWorkout.Name {
		t.Errorf("The names don't match %s, got %s", testWorkout.Name, workout.Name)
	}
}

func TestWorkoutService_GetByID_WrongUser(t *testing.T) {
	repo := newMockWorkoutRepository()
	service := NewWorkoutService(repo)

	testWorkout := createTestWorkout(t)

	err := service.Create(context.Background(), testWorkout)
	if err != nil {
		t.Errorf("Create failed %v", err)
	}

	workout, err := service.GetByID(context.Background(), testWorkout.ID, 2)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if workout != nil {
		t.Error("expected nil for wrong user, got workout")
	}
}

func TestWorkoutService_Update(t *testing.T) {
	repo := newMockWorkoutRepository()
	service := NewWorkoutService(repo)

	testWorkout := createTestWorkout(t)

	err := service.Create(context.Background(), testWorkout)
	if err != nil {
		t.Errorf("Create failed %v", err)
	}

	testWorkout.Name = "Legs"
	testWorkout.Notes = "Easy"

	err = service.Update(context.Background(), testWorkout)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	workout, err := service.GetByID(context.Background(), testWorkout.ID, 1)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if workout.Name != "Legs" {
		t.Errorf("expected name %s, got %s", testWorkout.Name, workout.Name)
	}
}

func TestWorkoutService_Delete(t *testing.T) {
	repo := newMockWorkoutRepository()
	service := NewWorkoutService(repo)

	testWorkout := createTestWorkout(t)

	err := service.Create(context.Background(), testWorkout)
	if err != nil {
		t.Errorf("Create failed %v", err)
	}

	err = service.Delete(context.Background(), testWorkout.ID, testWorkout.UserID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	workout, err := service.GetByID(context.Background(), testWorkout.ID, 1)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if workout != nil {
		t.Error("expected nil after delete, got workout")
	}
}

func TestWorkoutService_Delete_WrongUser(t *testing.T) {
	repo := newMockWorkoutRepository()
	service := NewWorkoutService(repo)

	testWorkout := createTestWorkout(t)

	err := service.Create(context.Background(), testWorkout)
	if err != nil {
		t.Errorf("Create failed %v", err)
	}

	err = service.Delete(context.Background(), testWorkout.ID, 2)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	workout, err := service.GetByID(context.Background(), testWorkout.ID, 1)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if workout == nil {
		t.Error("workout should still exist, but was deleted")
	}
}
