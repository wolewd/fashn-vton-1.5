package tryon

import (
	"archive/tar"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"

	"backend/utils"
)

func DownloadHandler(c echo.Context) error {
	root, _ := os.Getwd()

	jobsBaseDir := os.Getenv("JOB_DIR")
	if jobsBaseDir == "" {
		jobsBaseDir = filepath.Join(root, "jobs")
	}

	uuid := c.Param("uuid")
	targetPath := jobsBaseDir
	filename := "all_jobs.tar"

	if uuid != "" {
		targetPath = filepath.Join(jobsBaseDir, uuid)
		filename = fmt.Sprintf("job_%s.tar", uuid)
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil || len(entries) == 0 {
		return utils.JSONError(c, http.StatusNotFound, "Nothing to download", "Target path is empty or does not exist")
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
	c.Response().Header().Set(echo.HeaderContentType, "application/x-tar")
	c.Response().WriteHeader(http.StatusOK)

	tw := tar.NewWriter(c.Response().Writer)
	defer tw.Close()

	err = filepath.Walk(targetPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil { return err }
		if fi.IsDir() { return nil }

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil { return err }

		header.Name = strings.TrimPrefix(file, targetPath)

		if err := tw.WriteHeader(header); err != nil { return err }

		f, err := os.Open(file)
		if err != nil { return err }
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})

	return err
}
