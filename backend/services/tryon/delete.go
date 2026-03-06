package tryon

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"

	"backend/utils"
)

func DeleteHandler(c echo.Context) error {
	root, _ := os.Getwd()

	jobsBaseDir := os.Getenv("JOB_DIR")
	if jobsBaseDir == "" {
		jobsBaseDir = filepath.Join(root, "jobs")
	}

	uuid := c.Param("uuid")

	if uuid != "" {
		targetPath := filepath.Join(jobsBaseDir, uuid)

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			return utils.JSONError(c, http.StatusNotFound, "Job folder not found", uuid)
		}

		if err := os.RemoveAll(targetPath); err != nil {
			return utils.JSONError(c, http.StatusInternalServerError, "Failed to delete job folder", err.Error())
		}

		return utils.JSONSuccess(c, nil, fmt.Sprintf("Job %s deleted successfully", uuid))
	}

	entries, err := os.ReadDir(jobsBaseDir)
	if err != nil {
		return utils.JSONSuccess(c, map[string]int{"deleted_count": 0}, "Jobs directory is already empty or missing")
	}

	count := 0
	for _, entry := range entries {
		err := os.RemoveAll(filepath.Join(jobsBaseDir, entry.Name()))
		if err == nil {
			count++
		}
	}

	return utils.JSONSuccess(c, map[string]int{"deleted_count": count}, "All jobs cleared successfully")
}
