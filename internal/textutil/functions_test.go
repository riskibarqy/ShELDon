package textutil

import "testing"

func TestChoose(t *testing.T) {
	if got := Choose("value", "fallback"); got != "value" {
		t.Fatalf("expected value, got %q", got)
	}
	if got := Choose(" ", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}
}

func TestNormalizeCode(t *testing.T) {
	input := "```go\nfmt.Println(\"hi\")\n```"
	expected := "fmt.Println(\"hi\")\n"
	if got := NormalizeCode(input); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
	if got := NormalizeCode("fmt.Println(\"hi\")"); got != expected {
		t.Fatalf("expected newline appended, got %q", got)
	}
}

func TestExtractFunction(t *testing.T) {
	src := `
package demo

func target(a int) int {
	if a > 0 {
		return a
	}
	return -a
}
`
	got := ExtractFunction(src, "target")
	if got == "" {
		t.Fatalf("expected extraction result")
	}
	if ExtractFunction(src, "missing") != "" {
		t.Fatalf("expected missing function to return empty string")
	}
}
