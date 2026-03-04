package health

import (
	"net/http"
	"github.com/labstack/echo/v4"

	"backend/utils"
)

func Handler(c echo.Context) error {
	output, err := Check()
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Health check failed", map[string]string{
			"error": err.Error(),
			"output": string(output),
		})
	}

	return utils.JSONSuccess(c, map[string]string{
		"output": string(output),
	}, "Health check passed")
}
