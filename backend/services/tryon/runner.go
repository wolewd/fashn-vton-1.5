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

func NewJob(baseDir string, cfg Config) (*TryOnJob, error) {
	id := uuid.New().String()
	jobDir, _ := filepath.Abs(filepath.Join(baseDir, id)) // Use Absolute paths
	outputDir := filepath.Join(jobDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}

	return &TryOnJob{
		ID:        id,
		JobDir:    jobDir,
		OutputDir: outputDir,
		Config:    cfg,
	}, nil
}

func (job *TryOnJob) Run(ctx context.Context, personImg, garmentImg, category, photoType string, timesteps int) ([]byte, error) {
	// The Python script is strict: category must be 'tops', 'bottoms', or 'one-pieces'
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
