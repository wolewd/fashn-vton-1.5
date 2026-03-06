package health

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/cli"
)

func Check() ([]byte, error) {
	currentDir, _ := os.Getwd()

	// ADDITION: Resolve CLI Paths from .env to match the rest of the system
	venvPython := os.Getenv("PYTHON_PATH")
	if venvPython == "" {
		venvPython = filepath.Join(currentDir, ".venv", "bin", "python")
	}

	scriptPath := os.Getenv("SCRIPT_PATH")
	if scriptPath == "" {
		scriptPath = filepath.Join(currentDir, "examples", "basic_inference.py")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err := cli.RunPythonCLI(ctx, venvPython, scriptPath, "--help")

	outputStr := string(output)
	// Check if the output looks like the expected help text
	if strings.Contains(outputStr, "usage:") || strings.Contains(outputStr, "FASHN VTON") {
		return output, nil
	}

	return output, err
}
