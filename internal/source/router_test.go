package source

import (
	"context"
	"os"
	"testing"
)

type testConfig struct {
	cmd  string
	args []string
}

func (tc testConfig) GetCommand() string { return tc.cmd }
func (tc testConfig) GetArgs() []string  { return tc.args }

func TestNewScannerCommand(t *testing.T) {
	cfg := testConfig{cmd: "echo test"}
	s, cancel, err := NewScanner(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil || cancel == nil {
		t.Fatalf("expected scanner and cancel func")
	}
	defer func() { _ = cancel() }()
}

func TestNewScannerFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil && !os.IsNotExist(err) {
			t.Fatalf("failed to remove temp file: %v", err)
		}
	}()
	if _, err := tmpFile.WriteString("line1\nline2\n"); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	cfg := testConfig{args: []string{tmpFile.Name()}}
	s, cancel, err := NewScanner(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil || cancel == nil {
		t.Fatalf("expected scanner and cancel func")
	}
	defer func() { _ = cancel() }()
}

func TestNewScannerPipe(t *testing.T) {
	cfg := testConfig{}

	s, cancel, err := NewScanner(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil || cancel == nil {
		t.Fatalf("expected scanner and cancel func")
	}
	defer func() { _ = cancel() }()
}
