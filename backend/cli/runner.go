package cli

import (
	"context"
	"os/exec"
)

func RunPythonCLI(ctx context.Context, pythonPath string, scriptPath string, args ...string) ([]byte, error) {
	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.CommandContext(ctx, pythonPath, cmdArgs...)
	return cmd.CombinedOutput()
}
