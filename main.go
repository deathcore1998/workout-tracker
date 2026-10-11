package main

import (
	"context"
	"fmt"
	"log"

	"github.com/deathcore1998/workout-tracker/internal/config"
	"github.com/deathcore1998/workout-tracker/internal/database"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	pool, err := database.Connect(ctx, cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer pool.Close()
	fmt.Printf("Server will run on port %s\n", cfg.ServerPort)
}
