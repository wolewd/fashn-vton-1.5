package tryon

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"backend/utils"
)

// Handler handles the /api/tryon POST request
func Handler(c echo.Context) error {
	root, _ := os.Getwd()
	projectRoot := filepath.Dir(root)

	cfg := Config{
		PythonPath: filepath.Join(projectRoot, ".venv", "bin", "python"),
		ScriptPath: filepath.Join(projectRoot, "examples", "basic_inference.py"),
		WeightsDir: os.Getenv("WEIGHT_DIR"),
	}

	// 2. Validate Inputs
	personFile, err := c.FormFile("person_image")
	if err != nil { return utils.JSONError(c, http.StatusBadRequest, "Missing person_image", err.Error()) }

	garmentFile, err := c.FormFile("garment_image")
	if err != nil { return utils.JSONError(c, http.StatusBadRequest, "Missing garment_image", err.Error()) }

	category := c.FormValue("category") // Must be 'tops', 'bottoms', or 'one-pieces'

	// 3. Initialize Job
	job, err := NewJob("jobs", cfg)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, "Job creation failed", err.Error())
	}

	// 4. Save files with Absolute Paths
	personPath := filepath.Join(job.JobDir, "person.jpeg")
	garmentPath := filepath.Join(job.JobDir, "garment.jpeg")

	if err := saveUploadedFile(personFile, personPath); err != nil { /* handle error */ }
	if err := saveUploadedFile(garmentFile, garmentPath); err != nil { /* handle error */ }

	// 5. Execute with longer timeout (ML takes time)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	output, err := job.Run(ctx, personPath, garmentPath, category, "model", 30)
	if err != nil {
		// IMPORTANT: Capture 'output' here because it contains the Python Traceback
		return utils.JSONError(c, http.StatusInternalServerError, "CLI Execution Error", string(output))
	}

	return utils.JSONSuccess(c, map[string]string{
		"job_id":     job.ID,
		"result_dir": job.OutputDir,
	}, "Processing complete")
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
