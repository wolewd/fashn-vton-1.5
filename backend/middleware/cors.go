package middleware

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewCORSConfig() echo.MiddlewareFunc {
	origins := os.Getenv("ALLOWED_ORIGINS")

	allowList := []string{"*"}
	if origins != "" {
		allowList = strings.Split(origins, ",")
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowList,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			"X-App-ID",
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.OPTIONS,
		},
	})
}
