package wizard

import "testing"

func TestFormatTomlNoIndent_RemovesIndentation(t *testing.T) {
	input := "[agents]\n  [agents.codex]\n    enabled = true\n"
	want := "[agents]\n[agents.codex]\nenabled = true\n"

	got, err := formatTomlNoIndent(input)
	if err != nil {
		t.Fatalf("formatTomlNoIndent error: %v", err)
	}
	if got != want {
		t.Fatalf("formatTomlNoIndent output mismatch:\n%s\n---\n%s", got, want)
	}
}

func TestFormatTomlNoIndent_PreservesMultilineContent(t *testing.T) {
	input := "[section]\n  key = \"\"\"\n  line one\n  line two\n  \"\"\"\n  other = 1\n"
	want := "[section]\nkey = \"\"\"\n  line one\n  line two\n  \"\"\"\nother = 1\n"

	got, err := formatTomlNoIndent(input)
	if err != nil {
		t.Fatalf("formatTomlNoIndent error: %v", err)
	}
	if got != want {
		t.Fatalf("formatTomlNoIndent output mismatch:\n%s\n---\n%s", got, want)
	}
}
