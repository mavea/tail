package app

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
)

func TestNewApp(t *testing.T) {
	var wg sync.WaitGroup
	ctx := context.Background()

	app, appCancel, err := NewApp(os.Stdout.Fd(), ctx, &wg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil || appCancel == nil {
		t.Fatal("expected app and cancel func")
	}

	if err := appCancel(); err != nil {
		t.Fatalf("cancel failed: %v", err)
	}
}

func TestSetDefaultStd(t *testing.T) {
	var wg sync.WaitGroup
	app, _, _ := NewApp(os.Stdout.Fd(), context.Background(), &wg)

	app.SetDefaultStd(os.Stdin, os.Stdout, os.Stderr)
	if app.stdin != os.Stdin {
		t.Fatal("stdin not set")
	}
}

func TestDetectAndSetLanguage(t *testing.T) {
	var wg sync.WaitGroup
	app, _, _ := NewApp(os.Stdout.Fd(), context.Background(), &wg)

	if err := app.DetectAndSetLanguage("en_US.UTF-8"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.lang == nil {
		t.Fatal("expected language to be set")
	}

	if err := app.DetectAndSetLanguage(""); err != nil {
		t.Fatalf("unexpected error for empty lang: %v", err)
	}
}

func TestGetLang(t *testing.T) {
	var wg sync.WaitGroup
	app, _, _ := NewApp(os.Stdout.Fd(), context.Background(), &wg)

	if err := app.DetectAndSetLanguage("en_US.UTF-8"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lang := app.GetLang()
	if lang == nil {
		t.Fatal("expected language to be returned")
	}
}

func TestMergeWindowSizeWithConsole(t *testing.T) {
	var wg sync.WaitGroup
	app, _, _ := NewApp(os.Stdout.Fd(), context.Background(), &wg)

	mockCfg := &mockAppCfg{maxLineCount: 5, maxCharsPerLine: 50}

	result, err := app.mergeWindowSizeWithConsole(mockCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.GetMaxLineCount() != 5 {
		t.Fatalf("expected maxLineCount to stay 5, got %d", result.GetMaxLineCount())
	}
	if result.GetMaxCharsPerLine() != 50 {
		t.Fatalf("expected maxCharsPerLine to stay 50, got %d", result.GetMaxCharsPerLine())
	}

	mockCfg = &mockAppCfg{maxLineCount: 0, maxCharsPerLine: 0}
	result, err = app.mergeWindowSizeWithConsole(mockCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.GetMaxLineCount() < 1 {
		t.Fatalf("expected merged maxLineCount >= 1, got %d", result.GetMaxLineCount())
	}
	if result.GetMaxCharsPerLine() < 1 {
		t.Fatalf("expected merged maxCharsPerLine >= 1, got %d", result.GetMaxCharsPerLine())
	}
}

func TestApplyConfig(t *testing.T) {
	var wg sync.WaitGroup
	a, _, _ := NewApp(os.Stdout.Fd(), context.Background(), &wg)

	cfg := &mockAppCfg{maxLineCount: 5, maxCharsPerLine: 50}
	if err := a.ApplyConfig(cfg); err != nil {
		t.Fatalf("ApplyConfig failed: %v", err)
	}

	cfg.validateErr = errors.New("bad cfg")
	if err := a.ApplyConfig(cfg); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRunWithoutConfig(t *testing.T) {
	var wg sync.WaitGroup
	a, _, _ := NewApp(os.Stdout.Fd(), context.Background(), &wg)
	a.SetDefaultStd(os.Stdin, os.Stdout, os.Stderr)

	if err := a.Run(); err == nil {
		t.Fatal("expected error when config is not applied")
	}
}

func TestRunWithHelpConfig(t *testing.T) {
	var wg sync.WaitGroup
	a, _, _ := NewApp(os.Stdout.Fd(), context.Background(), &wg)
	a.SetDefaultStd(os.Stdin, os.Stdout, os.Stderr)
	if err := a.DetectAndSetLanguage("en_US.UTF-8"); err != nil {
		t.Fatalf("DetectAndSetLanguage failed: %v", err)
	}

	cfg := &mockAppCfg{maxLineCount: 5, maxCharsPerLine: 50, help: true}
	if err := a.ApplyConfig(cfg); err != nil {
		t.Fatalf("ApplyConfig failed: %v", err)
	}
	if err := a.Run(); err != nil {
		t.Fatalf("Run failed: %v", err)
	}
}

func TestCancelCollectsErrors(t *testing.T) {
	a := &App{}
	a.cancels = []func() error{
		func() error { return nil },
		func() error { return errors.New("one") },
		func() error { return errors.New("two") },
	}

	err := a.cancel()
	if err == nil {
		t.Fatal("expected joined cancel error")
	}
}

type mockAppCfg struct {
	maxLineCount    int
	maxCharsPerLine int
	validateErr     error
	help            bool
	version         bool
}

func (m *mockAppCfg) GetMaxCharsPerLine() int   { return m.maxCharsPerLine }
func (m *mockAppCfg) SetMaxCharsPerLine(n int)  { m.maxCharsPerLine = n }
func (m *mockAppCfg) GetMaxLineCount() int      { return m.maxLineCount }
func (m *mockAppCfg) SetMaxLineCount(n int)     { m.maxLineCount = n }
func (m *mockAppCfg) GetMaxBufferLines() uint64 { return 1000 }
func (m *mockAppCfg) GetProcessName() string    { return "" }
func (m *mockAppCfg) SetProcessName(string)     {}
func (m *mockAppCfg) GetProcessIcon() string    { return "" }
func (m *mockAppCfg) SetProcessIcon(string)     {}
func (m *mockAppCfg) GetOutputTemplate() string { return "none" }
func (m *mockAppCfg) GetIndicator() string      { return "none" }
func (m *mockAppCfg) GetOutputMode() string     { return "direct" }
func (m *mockAppCfg) IsHelp() bool              { return m.help }
func (m *mockAppCfg) IsVersion() bool           { return m.version }
func (m *mockAppCfg) GetCommand() string        { return "" }
func (m *mockAppCfg) GetArgs() []string         { return nil }
func (m *mockAppCfg) Validate() error           { return m.validateErr }
func (m *mockAppCfg) IsCSIEnabled() bool        { return true }
func (m *mockAppCfg) IsFullOutput() bool        { return false }
