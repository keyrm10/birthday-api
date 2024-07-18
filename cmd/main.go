package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/keyrm10/birthday-api/config"
	"github.com/keyrm10/birthday-api/internal/application/service"
	"github.com/keyrm10/birthday-api/internal/infrastructure/persistence"
	"github.com/keyrm10/birthday-api/internal/interfaces/api"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("application failed to start: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	db, err := persistence.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	if err := persistence.RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	userRepo := persistence.NewPgUserRepository(db)
	userService := service.NewUserService(userRepo)
	handler := api.NewHandler(userService)

	r := setupRouter(handler)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return gracefulShutdown(srv)
}

func setupRouter(handler *api.Handler) *gin.Engine {
	r := gin.Default()
	r.PUT("/hello/:username", handler.SaveUser)
	r.GET("/hello/:username", handler.GetUserBirthday)
	return r
}

func gracefulShutdown(srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exiting")
	return nil
}
