package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/deathcore1998/workout-tracker/internal/config"
	"github.com/deathcore1998/workout-tracker/internal/database"
	"github.com/deathcore1998/workout-tracker/internal/handlers"
	"github.com/deathcore1998/workout-tracker/internal/logger"
	"github.com/deathcore1998/workout-tracker/internal/middleware"
	"github.com/deathcore1998/workout-tracker/internal/repositories"
	"github.com/deathcore1998/workout-tracker/internal/services"
)

func main() {

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	logger.Setup(env)
	slog.Info("starting workout-tracker API", "env", env)

	ctx := context.Background()
	cfg := config.Load()
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer pool.Close()

	userRepo := repositories.NewUserRepository(pool)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, cfg.JWTSecret)

	workoutRepo := repositories.NewWorkoutRepository(pool)
	workoutService := services.NewWorkoutService(workoutRepo)
	workoutHandler := handlers.NewWorkoutHandler(workoutService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	authMiddleware := middleware.Auth(cfg.JWTSecret)
	mux.Handle("GET /profile", authMiddleware(http.HandlerFunc(userHandler.GetProfile)))

	mux.Handle("POST /workouts", authMiddleware(http.HandlerFunc(workoutHandler.Create)))
	mux.Handle("GET /workouts", authMiddleware(http.HandlerFunc(workoutHandler.List)))
	mux.Handle("GET /workouts/{id}", authMiddleware(http.HandlerFunc(workoutHandler.Get)))
	mux.Handle("PUT /workouts/{id}", authMiddleware(http.HandlerFunc(workoutHandler.Update)))
	mux.Handle("DELETE /workouts/{id}", authMiddleware(http.HandlerFunc(workoutHandler.Delete)))

	fmt.Printf("Server running on http://localhost: %s\n", cfg.ServerPort)

	loggedMux := middleware.Logger(mux)
	if err := http.ListenAndServe(":"+cfg.ServerPort, loggedMux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
