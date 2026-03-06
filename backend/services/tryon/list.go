package tryon

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"backend/utils"
)

type JobInfo struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
}

func ListHandler(c echo.Context) error {
	root, _ := os.Getwd()
	projectRoot := filepath.Dir(root)
	jobsBaseDir := filepath.Join(projectRoot, "jobs")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 { page = 1 }
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 { limit = 10 }

	entries, err := os.ReadDir(jobsBaseDir)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Failed to read jobs directory", err.Error())
	}

	var jobs []JobInfo
	for _, entry := range entries {
		if entry.IsDir() {
			info, _ := entry.Info()
			jobs = append(jobs, JobInfo{
				ID:        entry.Name(),
				CreatedAt: info.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
	}

	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].CreatedAt > jobs[j].CreatedAt
	})

	total := len(jobs)
	start := (page - 1) * limit
	end := start + limit

	if start > total { start = total }
	if end > total { end = total }

	paginatedJobs := jobs[start:end]

	return utils.JSONSuccess(c, map[string]interface{}{
		"jobs": paginatedJobs,
		"meta": map[string]interface{}{
			"limit":       limit,
			"page":        page,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	}, "Jobs retrieved successfully")
}
