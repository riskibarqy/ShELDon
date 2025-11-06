package analysis

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/riskiramdan/ShELDon/internal/system"
)

func TestImplScanner(t *testing.T) {
	dir := t.TempDir()
	source := `
package handlers

func register(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {})
}

func helper() {}
`
	file := filepath.Join(dir, "handler.go")
	if err := os.WriteFile(file, []byte(source), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	scanner := NewImplScanner(system.NewOSFileManager(strings.NewReader("")))
	out, err := scanner.Scan(dir)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}

	if !strings.Contains(out, "handler.go:5: r.GET") {
		t.Fatalf("expected handler snippet in output, got %q", out)
	}
}
