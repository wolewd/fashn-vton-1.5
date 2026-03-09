//go:build !windows

package cli

import (
	"context"
	"os/exec"
	"syscall"
)

func RunPythonCLI(ctx context.Context, pythonPath string, scriptPath string, args ...string) ([]byte, error) {
	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.CommandContext(ctx, pythonPath, cmdArgs...)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Cancel = func() error {
		if cmd.Process != nil {
			return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil
	}

	return cmd.CombinedOutput()
}