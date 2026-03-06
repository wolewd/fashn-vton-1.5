package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	appID := os.Getenv("APP_ID")

	// Create claims WITHOUT the "exp" field
	claims := jwt.MapClaims{
		"sub":     "vton-app",
		"iat":     time.Now().Unix(),
		// "exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	fmt.Println("Token Generated:")
	fmt.Printf("X-App-ID:      %s\n", appID)
	fmt.Printf("Authorization: Bearer %s\n", tokenString)
}
