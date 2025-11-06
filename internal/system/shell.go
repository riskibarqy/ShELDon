package system

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Shell abstracts running shell commands.
type Shell interface {
	Run(command string) (string, error)
}

// BashShell executes commands using `bash -lc`.
type BashShell struct{}

// Run executes the provided command string and returns stdout.
func (BashShell) Run(command string) (string, error) {
	if command == "" {
		return "", nil
	}
	cmd := exec.Command("bash", "-lc", command)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("shell command failed: %w\n%s", err, stderr.String())
	}
	return string(out), nil
}
