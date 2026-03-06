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
	projectRoot := filepath.Dir(root)
	jobsBaseDir := filepath.Join(projectRoot, "jobs")

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

	files, err := os.ReadDir(jobsBaseDir)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to read jobs directory", err.Error())
	}

	count := 0
	for _, f := range files {
		err := os.RemoveAll(filepath.Join(jobsBaseDir, f.Name()))
		if err == nil {
			count++
		}
	}

	return utils.JSONSuccess(c, map[string]int{"deleted_count": count}, "All jobs cleared successfully")
}
