package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"backend/middleware"
	"backend/services/health"
	"backend/services/tryon"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.Use(middleware.NewCORSConfig())
	e.Use(middleware.Logger)
	e.Use(middleware.NewRateLimiter(middleware.RateLimiterConfig{
		Requests: 5,
		Window:   time.Second,
	}))

	api := e.Group("/api")
	api.Use(middleware.AppIDAuth)
	api.Use(middleware.NewJWTAuth())
	{
		api.GET("/health", health.Handler)
		api.POST("/tryon", tryon.Handler)
	}

	go func() {
		if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("Shutting down the server due to error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited cleanly")
}
