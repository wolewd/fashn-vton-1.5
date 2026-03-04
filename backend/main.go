package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"backend/middleware"
	"backend/services/health"
	"backend/services/tryon"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using defaults")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.Use(middleware.Logger)
	e.Use(middleware.NewRateLimiter(middleware.RateLimiterConfig{
		Requests: 5,
		Window:   time.Second,
	}))

	e.GET("/api/health", health.Handler)
	e.POST("/api/tryon", tryon.Handler)
	e.Logger.Fatal(e.Start(":" + port))
}
