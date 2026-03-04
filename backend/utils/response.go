package utils

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func JSONSuccess(c echo.Context, data any, message string) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":  "success",
		"code":    200,
		"message": message,
		"data":    data,
	})
}

func JSONError(c echo.Context, code int, message string, err any) error {
	return c.JSON(code, map[string]any{
		"status":  "error",
		"code":    code,
		"message": message,
		"error":   err,
	})
}
