//go:build unix || linux || darwin

package builder

import "os/exec"

func createExec(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
