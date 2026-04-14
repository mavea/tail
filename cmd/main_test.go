package main

import (
	"os"
	"testing"

	"tail/internal/formatter"
)

func TestExitCodeFromError(t *testing.T) {
	if code := exitCodeFromError(nil); code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
	if code := exitCodeFromError(formatter.ExitError{Code: 7}); code != 7 {
		t.Fatalf("expected 7, got %d", code)
	}
	if code := exitCodeFromError(os.ErrInvalid); code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}

func TestExecuteHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"tail", "--help"}
	if err := execute(); err != nil {
		t.Fatalf("execute --help failed: %v", err)
	}
}

func TestExecuteInvalidFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"tail", "--unknown-flag"}
	if err := execute(); err == nil {
		t.Fatal("expected execute to fail on unknown flag")
	}
}
