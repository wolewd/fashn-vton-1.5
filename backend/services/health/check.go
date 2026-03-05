package health

import (
	"time"
	"os"
	"path/filepath"
	"context"
	"backend/cli"
	"strings"
)

func Check() ([]byte, error) {
    currentDir, _ := os.Getwd()
    projectRoot := filepath.Dir(currentDir)
    venvPython := filepath.Join(projectRoot, ".venv", "bin", "python")
    scriptPath := filepath.Join(projectRoot, "examples", "basic_inference.py")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    output, err := cli.RunPythonCLI(ctx, venvPython, scriptPath, "--help")

    outputStr := string(output)
    if strings.Contains(outputStr, "usage:") || strings.Contains(outputStr, "FASHN VTON") {
        return output, nil
    }

    return output, err
}
