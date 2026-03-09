package tryon

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"

	"backend/utils"
)

func ResultHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return utils.JSONError(c, http.StatusBadRequest, "UUID is required", nil)
	}

	root, _ := os.Getwd()
	jobsBaseDir := os.Getenv("JOB_DIR")
	if jobsBaseDir == "" {
		jobsBaseDir = filepath.Join(root, "jobs")
	}

	filePath := filepath.Join(jobsBaseDir, uuid, "output", "output_00.png")

	imageData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.JSONError(c, http.StatusNotFound, "Image not found or still processing", err.Error())
		}
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to read image file", err.Error())
	}

	base64Encoded := base64.StdEncoding.EncodeToString(imageData)

	mimeType := http.DetectContentType(imageData)

	return utils.JSONSuccess(c, map[string]interface{}{
		"uuid":      uuid,
		"mime_type": mimeType,
		"image":     base64Encoded,
	}, "Image retrieved successfully")
}
