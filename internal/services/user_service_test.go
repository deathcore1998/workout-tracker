package services

import (
	"context"
	"testing"

	"github.com/deathcore1998/workout-tracker/internal/models"
)

type mockUserRepository struct {
	users map[string]*models.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{users: make(map[string]*models.User)}
}

func (mock *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = len(mock.users) + 1
	mock.users[user.Email] = user
	return nil
}

func (mock *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := mock.users[email]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (mock *mockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	for _, user := range mock.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func createTestUser(t *testing.T) *models.User {
	t.Helper()

	return &models.User{
		Email:    "test@mail.com",
		Password: "password123",
		Name:     "Test User",
	}
}

func TestUserService_Register(t *testing.T) {
	repo := newMockUserRepository()
	service := NewUserService(repo)

	testUser := createTestUser(t)
	user, err := service.Register(context.Background(), testUser.Email, testUser.Password, testUser.Name)
	if err != nil {
		t.Fatalf("Register failed %v", err)
	}

	if user.ID == 0 {
		t.Error("expected user ID, got 0")
	}

	if user.Email != testUser.Email {
		t.Errorf("The email does not match %s, got %s", testUser.Email, user.Email)
	}

	if user.Password == testUser.Password {
		t.Error("password should be hashed, not plain text")
	}

	if user.Name != testUser.Name {
		t.Errorf("The name does not match %s, got %s", testUser.Name, user.Name)
	}
}

func TestUserService_Login(t *testing.T) {
	repo := newMockUserRepository()
	service := NewUserService(repo)

	testUser := createTestUser(t)
	_, err := service.Register(context.Background(), testUser.Email, testUser.Password, testUser.Name)
	if err != nil {
		t.Fatalf("Register failed %v", err)
	}

	user, err := service.Login(context.Background(), testUser.Email, testUser.Password)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if user.Email != testUser.Email {
		t.Errorf("expected email %s, got %s", testUser.Email, user.Email)
	}
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	repo := newMockUserRepository()
	service := NewUserService(repo)

	testUser := createTestUser(t)
	_, err := service.Register(context.Background(), testUser.Email, testUser.Password, testUser.Name)
	if err != nil {
		t.Fatalf("Register failed %v", err)
	}

	_, err = service.Login(context.Background(), testUser.Email, "wrongpassword")
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	repo := newMockUserRepository()
	service := NewUserService(repo)

	_, err := service.Login(context.Background(), "nonexistent@mail.com", "password123")
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

func TestUserService_GetByID(t *testing.T) {
	repo := newMockUserRepository()
	service := NewUserService(repo)

	testUser := createTestUser(t)
	created, err := service.Register(context.Background(), testUser.Email, testUser.Password, testUser.Name)
	if err != nil {
		t.Fatalf("Register failed %v", err)
	}

	user, err := service.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if user.Email != testUser.Email {
		t.Errorf("expected email %s, got %s", testUser.Email, user.Email)
	}
}
