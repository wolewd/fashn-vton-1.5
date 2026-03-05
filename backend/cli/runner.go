package cli

import (
	"context"
	"os/exec"
	"syscall"
)

func RunPythonCLI(ctx context.Context, pythonPath string, scriptPath string, args ...string) ([]byte, error) {
	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.CommandContext(ctx, pythonPath, cmdArgs...)

	cmd.Cancel = func() error {
		return cmd.Process.Signal(syscall.SIGKILL)
	}

	return cmd.CombinedOutput()
}
