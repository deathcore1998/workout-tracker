package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/deathcore1998/workout-tracker/internal/config"
	"github.com/deathcore1998/workout-tracker/internal/database"
	"github.com/deathcore1998/workout-tracker/internal/handlers"
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
	userHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", userHandler.Register)

	fmt.Printf("Server running on http://localhost: %s\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
