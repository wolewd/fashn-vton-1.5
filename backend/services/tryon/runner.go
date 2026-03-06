package tryon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"backend/cli"
)

type Config struct {
	PythonPath string
	ScriptPath string
	WeightsDir string
}

type TryOnJob struct {
	ID        string
	JobDir    string
	OutputDir string
	Config    Config
}

// getServicePaths centralized folder resolution
func getServicePaths() (string, string) {
	root, _ := os.Getwd()

	// 1. Resolve Job Directory
	jobDir := os.Getenv("JOB_DIR")
	if jobDir == "" {
		// Default to local jobs folder if environment variable is missing
		jobDir = filepath.Join(root, "jobs")
	}

	// 2. Resolve Weights Directory
	weightsDir := os.Getenv("WEIGHT_DIR")
	if weightsDir == "" {
		// Default to local weights folder if environment variable is missing
		weightsDir = filepath.Join(root, "weights")
	}

	return jobDir, weightsDir
}

func NewJob(cfg Config) (*TryOnJob, error) {
	jobBaseDir, _ := getServicePaths()
	id := uuid.New().String()

	// 3. Resolve Absolute Path
	// This works whether jobBaseDir is "../jobs" or "/home/wol/.../jobs"
	jobDir, err := filepath.Abs(filepath.Join(jobBaseDir, id))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	outputDir := filepath.Join(jobDir, "output")

	// 4. Create the job structure on disk
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create job directories: %w", err)
	}

	return &TryOnJob{
		ID:        id,
		JobDir:    jobDir,
		OutputDir: outputDir,
		Config:    cfg,
	}, nil
}

func (job *TryOnJob) Run(ctx context.Context, personImg, garmentImg, category, photoType string, timesteps int) ([]byte, error) {
	args := []string{
		"--weights-dir", job.Config.WeightsDir,
		"--person-image", personImg,
		"--garment-image", garmentImg,
		"--category", category,
		"--garment-photo-type", photoType,
		"--output-dir", job.OutputDir,
		"--num-timesteps", fmt.Sprint(timesteps),
		"--num-samples", "1",
	}

	return cli.RunPythonCLI(ctx, job.Config.PythonPath, job.Config.ScriptPath, args...)
}
