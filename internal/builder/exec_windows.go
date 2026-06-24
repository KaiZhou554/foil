//go:build windows

package builder

import "os/exec"

func createExec(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
