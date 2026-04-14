package formatter

import (
	"context"
	"testing"

	formatterWindow "tail/internal/formatter/window"
)

type mockOut struct {
	printCount int
	data       []string
	dataCount  int
	errorCount int
}

func (m *mockOut) GetDefaultStyle() (string, string, string) {
	return "", "", ""
}

func (m *mockOut) Print() error {
	m.printCount++
	return nil
}

func (m *mockOut) SetData(data []string) {
	m.data = data
	m.dataCount++
}

func (m *mockOut) ClearScreen() error {
	return nil
}

func (m *mockOut) Error(_ []string, _ []string) error {
	m.errorCount++
	return nil
}

func (m *mockOut) SetStatus(_ int) error {
	return nil
}

func (m *mockOut) PrintFull(_ []string) error {
	return nil
}

type mockCfg struct {
	mode            string
	maxLineCount    int
	maxCharsPerLine int
	maxBufferLines  uint64
	outputTemplate  string
	indicator       string
	processName     string
	processIcon     string
}

func (mc *mockCfg) GetOutputMode() string     { return mc.mode }
func (mc *mockCfg) GetMaxLineCount() int      { return mc.maxLineCount }
func (mc *mockCfg) GetMaxCharsPerLine() int   { return mc.maxCharsPerLine }
func (mc *mockCfg) GetMaxBufferLines() uint64 { return mc.maxBufferLines }
func (mc *mockCfg) GetOutputTemplate() string { return mc.outputTemplate }
func (mc *mockCfg) GetIndicator() string      { return mc.indicator }
func (mc *mockCfg) GetProcessName() string    { return mc.processName }
func (mc *mockCfg) GetProcessIcon() string    { return mc.processIcon }
func (mc *mockCfg) IsCSIEnabled() bool        { return true }
func (mc *mockCfg) IsFullOutput() bool        { return false }

func TestNewFormatterDirect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := &mockOut{}
	cfg := &mockCfg{mode: "direct", maxLineCount: 10, maxBufferLines: 1000}
	window := formatterWindow.NewWindow("", "")

	r, cancelFn, err := NewFormatter(ctx, out, cfg, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil || cancelFn == nil {
		t.Fatal("expected formatter and cancel func")
	}

	if err := r.Set("line1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if out.dataCount == 0 {
		t.Fatal("expected SetData to be called")
	}
	if out.printCount == 0 {
		t.Fatal("expected Print to be called in direct mode")
	}

	// Cancel will wait for ongoing goroutines and clean up resources
	_ = cancelFn()
}

func TestNewFormatterThread(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := &mockOut{}
	cfg := &mockCfg{mode: "thread", maxLineCount: 10, maxBufferLines: 1000}
	window := formatterWindow.NewWindow("", "")

	r, cancelFn, err := NewFormatter(ctx, out, cfg, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := r.Set("line1"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if out.printCount > 0 {
		t.Fatal("expected no immediate Print in thread mode")
	}

	_ = cancelFn()
}

func TestFormatterSetErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := &mockOut{}
	cfg := &mockCfg{mode: "direct", maxLineCount: 10, maxBufferLines: 1000}
	window := formatterWindow.NewWindow("", "")

	r, cancelFn, err := NewFormatter(ctx, out, cfg, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := r.SetErr(ErrTest()); err != nil {
		t.Fatalf("SetErr failed: %v", err)
	}

	if err := cancelFn(); err == nil {
		t.Fatal("expected cancel to return error when errors were set")
	}
}

func TestFormatterSetErrLineDoesNotFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := &mockOut{}
	cfg := &mockCfg{mode: "direct", maxLineCount: 10, maxBufferLines: 1000}
	window := formatterWindow.NewWindow("", "")

	r, cancelFn, err := NewFormatter(ctx, out, cfg, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := r.SetErrLine("warning"); err != nil {
		t.Fatalf("SetErrLine failed: %v", err)
	}

	if err := cancelFn(); err != nil {
		t.Fatalf("expected non-fatal cancel for stderr line, got: %v", err)
	}
	if out.errorCount != 0 {
		t.Fatalf("expected no error rendering for non-fatal stderr, got: %d", out.errorCount)
	}
}

func TestFormatterSetStatus(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := &mockOut{}
	cfg := &mockCfg{mode: "direct", maxLineCount: 10, maxBufferLines: 1000}
	window := formatterWindow.NewWindow("", "")

	r, cancelFn, err := NewFormatter(ctx, out, cfg, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := r.SetStatus(42); err != nil {
		t.Fatalf("SetStatus failed: %v", err)
	}

	cancelErr := cancelFn()
	if cancelErr == nil {
		t.Fatal("expected cancel to return error when status is non-zero")
	}
}

type testErr struct{}

func (testErr) Error() string { return "test error" }

func ErrTest() error {
	return testErr{}
}
