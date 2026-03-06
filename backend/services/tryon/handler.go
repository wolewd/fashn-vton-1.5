package tryon

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"backend/utils"
)

func Handler(c echo.Context) error {
	root, _ := os.Getwd()

	_, weightsDir := getServicePaths()

	timesteps, _ := strconv.Atoi(os.Getenv("NUM_TIMESTEP"))
	if timesteps <= 0 {
		timesteps = 30
	}

	// 1. Resolve CLI Paths: Use .env values (Absolute or Relative)
	// This will pick up your "../" locally and "/root/..." in Docker
	pythonPath := os.Getenv("PYTHON_PATH")
	if pythonPath == "" {
		// Only fallback to a local path if the .env variable is missing
		pythonPath = filepath.Join(root, ".venv", "bin", "python")
	}

	scriptPath := os.Getenv("SCRIPT_PATH")
	if scriptPath == "" {
		scriptPath = filepath.Join(root, "examples", "basic_inference.py")
	}

	cfg := Config{
		PythonPath: pythonPath,
		ScriptPath: scriptPath,
		WeightsDir: weightsDir,
	}

	// 2. Initialize Job
	job, err := NewJob(cfg)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Job creation failed", err.Error())
	}

	errWithID := func(code int, msg string, detail any) error {
		return c.JSON(code, map[string]any{
			"status":  "error",
			"code":    code,
			"message": msg,
			"error":   detail,
			"job_id":  job.ID,
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

	// 5. Execute with detailed error capture
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Minute)
	defer cancel()

	output, err := job.Run(ctx, personPath, garmentPath, category, photoType, timesteps)
	if err != nil {
		errorMessage := string(output)
		if errorMessage == "" {
			errorMessage = err.Error()
		}
		return errWithID(http.StatusInternalServerError, "CLI Execution Error", errorMessage)
	}

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
