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
}

// OSFileManager implements FileManager using the host filesystem.
type OSFileManager struct {
	Stdin io.Reader
}

// NewOSFileManager creates a FileManager backed by os package primitives.
func NewOSFileManager(stdin io.Reader) *OSFileManager {
	if stdin == nil {
		stdin = os.Stdin
	}
	return &OSFileManager{Stdin: stdin}
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

func (fm *OSFileManager) readFromStdin() (string, error) {
	if fm.Stdin == nil {
		return "", errors.New("stdin is not available")
	}
	bytes, err := io.ReadAll(fm.Stdin)
	return string(bytes), err
}
