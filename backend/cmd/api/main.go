package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qris-pos-backend/internal/infrastructure/config"
	"qris-pos-backend/internal/infrastructure/database"
	"qris-pos-backend/internal/interfaces/http/server"
	"qris-pos-backend/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.App.LogLevel)

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", "error", err)
	}
	defer database.Close(db)

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		appLogger.Fatal("Failed to run migrations", "error", err)
	}

	// Seed initial data
	if err := database.SeedData(db); err != nil {
		appLogger.Fatal("Failed to seed data", "error", err)
	}

	// Initialize HTTP server
	httpServer := server.NewServer(cfg, db, appLogger)

	// Start server in a goroutine
	go func() {
		appLogger.Info("Starting server", "port", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exited")
}