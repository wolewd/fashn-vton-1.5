package tryon

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"backend/utils"
)

// Handler handles the /api/tryon POST request
func Handler(c echo.Context) error {
	// Get uploaded files
	personFile, err := c.FormFile("person_image")
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "person_image is required", "")
	}
	garmentFile, err := c.FormFile("garment_image")
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "garment_image is required", "")
	}

	category := c.FormValue("category")
	if category == "" {
		return utils.JSONError(c, http.StatusBadRequest, "category is required", "")
	}

	// Optional garment photo type
	garmentPhotoType := c.FormValue("garment_photo_type")
	if garmentPhotoType == "" {
		garmentPhotoType = "model"
	}

	// Create a new job folder
	job, err := NewJob("jobs")
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to create job directory", err.Error())
	}

	// Save uploaded files
	personPath := filepath.Join(job.JobDir, "person.jpeg")
	garmentPath := filepath.Join(job.JobDir, "garment.jpeg")

	if err := saveUploadedFile(personFile, personPath); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to save person image", err.Error())
	}
	if err := saveUploadedFile(garmentFile, garmentPath); err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "failed to save garment image", err.Error())
	}

	// Load environment variables
	weightsDir := os.Getenv("WEIGHT_DIR")
	if weightsDir == "" {
		weightsDir = "./weights"
	}

	numTimesteps := 30
	if v := os.Getenv("NUM_TIMESTEP"); v != "" {
		if t, err := strconv.Atoi(v); err == nil {
			numTimesteps = t
		}
	}

	pythonPath := filepath.Join(".", ".venv", "bin", "python")
	scriptPath := filepath.Join(".", "examples", "basic_inference.py")

	// Run the Python CLI using the job runner
	output, err := job.Run(
		pythonPath,
		scriptPath,
		weightsDir,
		personPath,
		garmentPath,
		category,
		garmentPhotoType,
		numTimesteps,
	)
	if err != nil {
		fmt.Println("CLI Output: ", string(output))
		return utils.JSONError(c, http.StatusInternalServerError, "Try-on CLI failed", string(output))
	}

	// Return standardized JSON success response
	return utils.JSONSuccess(c, map[string]string{
		"job_id":     job.ID,
		"output_dir": job.OutputDir,
	}, "Try-on job completed")
}

// saveUploadedFile saves an uploaded *multipart.FileHeader to a destination path
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
