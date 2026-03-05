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
    // 1. Resolve Paths
    root, _ := os.Getwd()
    // If running from within 'backend', move up to project root
    projectRoot := filepath.Dir(root)

    // Ensure we use the absolute path for weights from the .env or default to root/weights
    weightsDir := os.Getenv("WEIGHT_DIR")
    if weightsDir == "" || weightsDir == "./weights" {
        weightsDir = filepath.Join(projectRoot, "weights")
    }

    cfg := Config{
        PythonPath: filepath.Join(projectRoot, ".venv", "bin", "python"),
        ScriptPath: filepath.Join(projectRoot, "examples", "basic_inference.py"),
        WeightsDir: weightsDir,
    }

    // 2. Validate Inputs
    personFile, err := c.FormFile("person_image")
    if err != nil { return utils.JSONError(c, http.StatusBadRequest, "Missing person_image", err.Error()) }

    garmentFile, err := c.FormFile("garment_image")
    if err != nil { return utils.JSONError(c, http.StatusBadRequest, "Missing garment_image", err.Error()) }

    category := c.FormValue("category")

    // 3. Initialize Job in root/jobs
    // We pass the absolute path to the root 'jobs' folder
    jobsBaseDir := filepath.Join(projectRoot, "jobs")
    job, err := NewJob(jobsBaseDir, cfg)
    if err != nil {
        return utils.JSONError(c, http.StatusInternalServerError, "Job creation failed", err.Error())
    }

    // 4. Save files with Absolute Paths
    personPath := filepath.Join(job.JobDir, "person.jpeg")
    garmentPath := filepath.Join(job.JobDir, "garment.jpeg")

    if err := saveUploadedFile(personFile, personPath); err != nil {
        return utils.JSONError(c, http.StatusInternalServerError, "Failed to save person image", err.Error())
    }
    if err := saveUploadedFile(garmentFile, garmentPath); err != nil {
        return utils.JSONError(c, http.StatusInternalServerError, "Failed to save garment image", err.Error())
    }

    // 5. Execute with 5-minute timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    // Pass 'category' from form; ensure garment_photo_type is also handled if needed
    photoType := c.FormValue("garment_photo_type")
    if photoType == "" { photoType = "model" }

    output, err := job.Run(ctx, personPath, garmentPath, category, photoType, 30)
    if err != nil {
        return utils.JSONError(c, http.StatusInternalServerError, "CLI Execution Error", string(output))
    }

    return utils.JSONSuccess(c, map[string]string{
        "job_id":     job.ID,
        "result_dir": job.OutputDir,
    }, "Processing complete")
}

func saveUploadedFile(fileHeader *multipart.FileHeader, dest string) error {
    src, err := fileHeader.Open()
    if err != nil { return err }
    defer src.Close()

    out, err := os.Create(dest)
    if err != nil { return err }
    defer out.Close()

    _, err = io.Copy(out, src)
    return err
}
