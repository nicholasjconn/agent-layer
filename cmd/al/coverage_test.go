package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootRunE_Version(t *testing.T) {

	cmd := newRootCmd()

	cmd.Version = "dev"

	if err := cmd.Flags().Set("version", "true"); err != nil {

		t.Fatalf("set flag: %v", err)

	}

	var out bytes.Buffer

	cmd.SetOut(&out)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE error: %v", err)
	}
	if strings.TrimSpace(out.String()) != "dev" {
		t.Fatalf("expected version output, got %q", out.String())
	}
}

func TestRootRunE_Help(t *testing.T) {
	cmd := newRootCmd()
	// Version default is false

	// We need to capture help output. cmd.Help() writes to Out.
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE error: %v", err)
	}
	// Help output should be present
	if !strings.Contains(out.String(), "Usage:") {
		t.Fatalf("expected help output")
	}
}
