package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Client abstracts interaction with the git binary.
type Client interface {
	Diff(args ...string) (string, error)
	Commit(message string) error
}

// CLIClient runs git commands via the local binary.
type CLIClient struct{}

// Diff returns the diff output for the provided arguments.
func (CLIClient) Diff(args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"diff"}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff %v: %w\n%s", args, err, stderr.String())
	}
	return string(out), nil
}

// Commit runs `git commit -m <message>`.
func (CLIClient) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &bytes.Buffer{}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w\n%s", err, stderr.String())
	}
	return nil
}
