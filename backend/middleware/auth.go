package middleware

import (
	"os"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"

	"backend/utils"
)

// AppIDAuth verify the X-App-ID header exists and matches our env
func AppIDAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		expectedAppID := os.Getenv("APP_ID")
		requestAppID := c.Request().Header.Get("X-App-ID")

		if expectedAppID == "" || requestAppID != expectedAppID {
			return utils.JSONError(c, http.StatusUnauthorized, "Invalid or missing App ID", nil)
		}
		return next(c)
	}
}

// NewJWTAuth creates a middleware that validates JWT tokens from .env
func NewJWTAuth() echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}

	config := echojwt.Config{
		SigningKey: []byte(secret),
		ContextKey: "user",
		ErrorHandler: func(c echo.Context, err error) error {
			return utils.JSONError(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
		},
	}

	return echojwt.WithConfig(config)
}

// GetUserID helper
func GetUserID(c echo.Context) string {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return ""
	}
	claims := user.Claims.(jwt.MapClaims)
	return claims["sub"].(string)
}
