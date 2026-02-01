package main

import (
	"context"
	"log"
	"time"

	"github.com/martiriera/discogs-spotify/internal/infrastructure/config"
	"github.com/martiriera/discogs-spotify/internal/infrastructure/container"
)

func main() {
	log.Println("Starting worker...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	c := container.NewContainer(cfg)

	// Create a context with timeout for the worker
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Example worker task - you can expand this based on your needs
	if err := runWorkerTask(ctx, c); err != nil {
		log.Fatalf("Worker task failed: %v", err)
	}

	log.Println("Worker completed successfully")
}

func runWorkerTask(ctx context.Context, c *container.Container) error {
	log.Println("Running worker task...")

	// Example: This could be batch processing, cleanup tasks, etc.
	// For now, it's a simple placeholder that demonstrates the pattern
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		log.Println("Worker task processing completed")
		return nil
	}
}
