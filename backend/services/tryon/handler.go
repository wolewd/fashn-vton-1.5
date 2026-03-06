package tryon

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"backend/utils"
	"github.com/labstack/echo/v4"
)

func Handler(c echo.Context) error {
	// 1. Resolve Paths
	root, _ := os.Getwd()
	projectRoot := filepath.Dir(root)

	weightsDir := os.Getenv("WEIGHT_DIR")
	if weightsDir == "" || weightsDir == "./weights" {
		weightsDir = filepath.Join(projectRoot, "weights")
	}

	cfg := Config{
		PythonPath: filepath.Join(projectRoot, ".venv", "bin", "python"),
		ScriptPath: filepath.Join(projectRoot, "examples", "basic_inference.py"),
		WeightsDir: weightsDir,
	}

	// 2. Initialize Job First (So we have a UUID for every response)
	jobsBaseDir := filepath.Join(projectRoot, "jobs")
	job, err := NewJob(jobsBaseDir, cfg)
	if err != nil {
		// No Job ID yet if this fails
		return utils.JSONError(c, http.StatusInternalServerError, "Job creation failed", err.Error())
	}

	// Helper to attach job_id to all errors from here on
	errWithID := func(code int, msg string, detail any) error {
		return c.JSON(code, map[string]any{
			"status":  "error",
			"code":    code,
			"message": msg,
			"error":   detail,
			"job_id":  job.ID, // Always return the ID
		})
	}

	// 3. Validate Inputs
	personFile, err := c.FormFile("person_image")
	if err != nil {
		return errWithID(http.StatusBadRequest, "Missing person_image", err.Error())
	}

	garmentFile, err := c.FormFile("garment_image")
	if err != nil {
		return errWithID(http.StatusBadRequest, "Missing garment_image", err.Error())
	}

	category := c.FormValue("category")
	photoType := c.FormValue("garment_photo_type")
	if photoType == "" {
		photoType = "model"
	}

	// 4. Save files
	personPath := filepath.Join(job.JobDir, "person.jpeg")
	garmentPath := filepath.Join(job.JobDir, "garment.jpeg")

	if err := saveUploadedFile(personFile, personPath); err != nil {
		return errWithID(http.StatusInternalServerError, "Failed to save person image", err.Error())
	}
	if err := saveUploadedFile(garmentFile, garmentPath); err != nil {
		return errWithID(http.StatusInternalServerError, "Failed to save garment image", err.Error())
	}

	// 5. Execute with Context Inheritance
	// We wrap the Request Context. If the user cancels the CURL, ctx is canceled.
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Minute)
	defer cancel()

	output, err := job.Run(ctx, personPath, garmentPath, category, photoType, 30)
	if err != nil {
		// Return Python error output + the job_id
		return errWithID(http.StatusInternalServerError, "CLI Execution Error", string(output))
	}

	// 6. Final Success
	return utils.JSONSuccess(c, map[string]string{
		"job_id":     job.ID,
		"result_dir": job.OutputDir,
	}, "Processing complete")
}

func saveUploadedFile(fileHeader *multipart.FileHeader, dest string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
