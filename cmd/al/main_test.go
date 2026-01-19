package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestMainVersion(t *testing.T) {
	var out bytes.Buffer
	if err := execute([]string{"al", "--version"}, &out, &out); err != nil {
		t.Fatalf("execute error: %v", err)
	}
	if !strings.Contains(out.String(), Version) {
		t.Fatalf("expected version output, got %q", out.String())
	}
}

func TestMainUnknownCommand(t *testing.T) {
	var out bytes.Buffer
	err := execute([]string{"al", "unknown"}, &out, &out)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunMainSuccess(t *testing.T) {
	var out bytes.Buffer
	called := false
	runMain([]string{"al", "--version"}, &out, &out, func(code int) {
		called = true
	})
	if called {
		t.Fatalf("unexpected exit")
	}
}

func TestRunMainError(t *testing.T) {
	var out bytes.Buffer
	code := 0
	runMain([]string{"al", "unknown"}, &out, &out, func(exitCode int) {
		code = exitCode
	})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(out.String(), "unknown command") {
		t.Fatalf("expected error output, got %q", out.String())
	}
}

func TestMainCallsExecute(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"al", "--version"}
	main()
}
