package commands

import "testing"

func TestApplyPrefix(t *testing.T) {
	msg := "feat: add feature\n\n* detail"
	with := applyPrefix(msg, "chore:")
	expected := "chore: feat: add feature\n\n* detail"
	if with != expected {
		t.Fatalf("expected %q, got %q", expected, with)
	}
	if applyPrefix(msg, "") != msg {
		t.Fatalf("expected unchanged message when prefix empty")
	}
}

func TestNormalizeCommitMessage(t *testing.T) {
	raw := `
Here is a concise message:
**feat: do thing**

* bullet
`
	want := "feat: do thing\n\n* bullet"
	if got := normalizeCommitMessage(raw); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestValidateConventionalCommitWithPrefix(t *testing.T) {
	cases := []struct {
		name    string
		message string
		valid   bool
	}{
		{"plain", "feat(parser): add validation", true},
		{"ticket prefix", "WIT-341 fix(service): handle edge case", true},
		{"multi prefix", "[WIP] NO-TICKET chore(ci): tweak workflow", true},
		{"invalid", "WIT-341 missing colon", false},
	}
	for _, tc := range cases {
		if got := validateConventionalCommit(tc.message); got != tc.valid {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.valid, got)
		}
	}
}
