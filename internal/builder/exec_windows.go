//go:build windows

package builder

import (
	"os/exec"
	"syscall"
)

func createExec(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// hiddenCmd creates a command that doesn't show a console window on Windows.
func hiddenCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
	return cmd
}
