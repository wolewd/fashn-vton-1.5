package tryon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"backend/cli"
)

type TryOnJob struct {
	ID        string
	JobDir    string
	OutputDir string
}

// NewJob creates a new job folder under baseDir
func NewJob(baseDir string) (*TryOnJob, error) {
	id := uuid.New().String()
	jobDir := filepath.Join(baseDir, id)
	outputDir := filepath.Join(jobDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}

	return &TryOnJob{
		ID:        id,
		JobDir:    jobDir,
		OutputDir: outputDir,
	}, nil
}

// Run executes the Python CLI for this job
func (job *TryOnJob) Run(
	pythonPath, scriptPath, weightsDir, personImage, garmentImage, category, garmentPhotoType string,
	numTimesteps int,
) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	args := []string{
		"--weights-dir", weightsDir,
		"--person-image", personImage,
		"--garment-image", garmentImage,
		"--category", category,
		"--garment-photo-type", garmentPhotoType,
		"--output-dir", job.OutputDir,
		"--num-timesteps", fmt.Sprint(numTimesteps),
	}

	return cli.RunPythonCLI(ctx, pythonPath, scriptPath, args...)
}
