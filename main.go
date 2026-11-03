package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/deathcore1998/workout-tracker/internal/config"
	"github.com/deathcore1998/workout-tracker/internal/database"
	"github.com/deathcore1998/workout-tracker/internal/handlers"
	"github.com/deathcore1998/workout-tracker/internal/middleware"
	"github.com/deathcore1998/workout-tracker/internal/repositories"
	"github.com/deathcore1998/workout-tracker/internal/services"
)

func main() {
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
	if err := http.ListenAndServe(":"+cfg.ServerPort, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
