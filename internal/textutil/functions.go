package textutil

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Choose returns the first non-empty value, defaulting to fallback.
func Choose(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

// NormalizeCode strips Markdown fences and ensures trailing newline.
func NormalizeCode(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		if len(lines) >= 2 {
			s = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	return s
}

// ExtractFunction attempts to extract the named Go function body.
func ExtractFunction(src, fn string) string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return ""
	}

	for _, decl := range file.Decls {
		fnDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fnDecl.Name.Name != fn {
			continue
		}

		start := fset.Position(fnDecl.Pos()).Offset
		end := fset.Position(fnDecl.End()).Offset
		if start < 0 || end < start || end > len(src) {
			return ""
		}
		return src[start:end]
	}
	return ""
}
