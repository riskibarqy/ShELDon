package system

import (
	"errors"
	"io"
	"os"
)

// FileManager abstracts file IO to keep command logic testable.
type FileManager interface {
	Read(path string) (string, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data string) error
	IsInteractive() bool
}

// OSFileManager implements FileManager using the host filesystem.
type OSFileManager struct {
	Stdin       io.Reader
	interactive bool
}

// NewOSFileManager creates a FileManager backed by os package primitives.
func NewOSFileManager(stdin io.Reader) *OSFileManager {
	if stdin == nil {
		stdin = os.Stdin
	}
	fm := &OSFileManager{Stdin: stdin}
	if stdin == os.Stdin {
		if info, err := os.Stdin.Stat(); err == nil {
			fm.interactive = info.Mode()&os.ModeCharDevice != 0
		}
	}
	return fm
}

// Read returns file content or reads from stdin when path is empty or "-".
func (fm *OSFileManager) Read(path string) (string, error) {
	switch path {
	case "", "-":
		return fm.readFromStdin()
	default:
		data, err := os.ReadFile(path)
		return string(data), err
	}
}

// ReadFile reads a file and returns the raw bytes.
func (fm *OSFileManager) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile persists data using 0644 permissions.
func (fm *OSFileManager) WriteFile(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0o644)
}

// IsInteractive reports whether stdin is attached to a terminal.
func (fm *OSFileManager) IsInteractive() bool {
	return fm.interactive
}

func (fm *OSFileManager) readFromStdin() (string, error) {
	if fm.Stdin == nil {
		return "", errors.New("stdin is not available")
	}
	bytes, err := io.ReadAll(fm.Stdin)
	return string(bytes), err
}
