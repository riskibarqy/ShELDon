package analysis

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/riskiramdan/ShELDon/internal/system"
)

// ImplScanner collects handler snippets from Go implementation files.
type ImplScanner struct {
	Files system.FileManager
}

// NewImplScanner creates a scanner with the provided file manager.
func NewImplScanner(files system.FileManager) ImplScanner {
	return ImplScanner{Files: files}
}

// Scan walks the directory capturing lines that likely describe handlers.
func (s ImplScanner) Scan(dir string) (string, error) {
	var builder strings.Builder
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		content, err := s.Files.Read(path)
		if err != nil {
			return err
		}
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if containsHandlerHint(line) {
				fmt.Fprintf(&builder, "%s:%d: %s\n", path, i+1, line)
			}
		}
		return nil
	})
	return builder.String(), err
}

func containsHandlerHint(line string) bool {
	return strings.Contains(line, "gin.") ||
		strings.Contains(line, "echo.") ||
		strings.Contains(line, "http.") ||
		strings.Contains(line, "fiber.")
}
