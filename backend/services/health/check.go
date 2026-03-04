package health

import (
	"time"
	"os"
	"path/filepath"
	"context"
	"backend/cli"
)

func Check() ([]byte, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	projectRoot := filepath.Dir(currentDir)
	venvPython := filepath.Join(projectRoot, ".venv", "bin", "python")
	scriptPath := filepath.Join(projectRoot, "examples", "basic_inference.py")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return cli.RunPythonCLI(ctx, venvPython, scriptPath, "--help")
}
