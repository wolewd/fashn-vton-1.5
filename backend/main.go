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

	host := os.Getenv("SERVER_HOST")
    if host == "" {
        host = "0.0.0.0"
    }

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)

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
		api.GET("/tryon/:uuid", tryon.ResultHandler)

		api.GET("/jobs", tryon.ListHandler)
        api.DELETE("/jobs", tryon.DeleteHandler)
        api.DELETE("/jobs/:uuid", tryon.DeleteHandler)

        api.GET("/downloads", tryon.DownloadHandler)
        api.GET("/downloads/:uuid", tryon.DownloadHandler)
	}

	go func() {
		if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
