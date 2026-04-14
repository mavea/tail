package config

import (
	"tail/internal/formatter"
	"testing"

	sinkConsoleIndicator "tail/internal/sink/console/indicator"
	sinkConsoleTemplate "tail/internal/sink/console/template"
)

func TestCfgGettersAndSetters(t *testing.T) {
	cfg := Cfg{
		maxLineCount:    10,
		maxCharsPerLine: 100,
		maxBufferLines:  1000,
		processName:     "test-proc",
		processIcon:     "🔧",
		outputMode:      formatter.NewOutputMode(),
		outputTemplate:  sinkConsoleTemplate.NewTemplateType(),
		indicator:       sinkConsoleIndicator.NewIndicatorType(),
		help:            false,
		version:         false,
		command:         "cmd",
		args:            []string{"arg1"},
	}

	if cfg.GetMaxLineCount() != 10 {
		t.Fatalf("unexpected maxLineCount: %d", cfg.GetMaxLineCount())
	}
	cfg.SetMaxLineCount(20)
	if cfg.GetMaxLineCount() != 20 {
		t.Fatalf("unexpected maxLineCount after set: %d", cfg.GetMaxLineCount())
	}

	if cfg.GetMaxCharsPerLine() != 100 {
		t.Fatalf("unexpected maxCharsPerLine: %d", cfg.GetMaxCharsPerLine())
	}
	cfg.SetMaxCharsPerLine(200)
	if cfg.GetMaxCharsPerLine() != 200 {
		t.Fatalf("unexpected maxCharsPerLine after set: %d", cfg.GetMaxCharsPerLine())
	}

	if cfg.GetMaxBufferLines() != 1000 {
		t.Fatalf("unexpected maxBufferLines: %d", cfg.GetMaxBufferLines())
	}

	if cfg.GetProcessName() != "test-proc" {
		t.Fatalf("unexpected processName: %s", cfg.GetProcessName())
	}

	if cfg.GetProcessIcon() != "🔧" {
		t.Fatalf("unexpected processIcon: %s", cfg.GetProcessIcon())
	}

	if !cfg.IsHelp() != true {
		t.Fatalf("unexpected help flag")
	}
	if !cfg.IsVersion() != true {
		t.Fatalf("unexpected version flag")
	}

	if cfg.GetCommand() != "cmd" {
		t.Fatalf("unexpected command: %s", cfg.GetCommand())
	}

	if len(cfg.GetArgs()) != 1 || cfg.GetArgs()[0] != "arg1" {
		t.Fatalf("unexpected args: %v", cfg.GetArgs())
	}

	if cfg.IsCSIEnabled() {
		t.Fatal("expected CSI to be not enabled")
	}
}

func TestCfgValidate(t *testing.T) {
	cfg := Cfg{
		maxLineCount:    10,
		maxCharsPerLine: 100,
		maxBufferLines:  1000,
		outputMode:      formatter.NewOutputMode(),
		outputTemplate:  sinkConsoleTemplate.NewTemplateType(),
		indicator:       sinkConsoleIndicator.NewIndicatorType(),
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid config failed validation: %v", err)
	}

	cfg.maxLineCount = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero maxLineCount")
	}

	cfg.maxLineCount = 251
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for maxLineCount > 250")
	}

	cfg.maxLineCount = 10
	cfg.maxCharsPerLine = -1
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative maxCharsPerLine")
	}

	cfg.maxCharsPerLine = 1001
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for maxCharsPerLine > 1000")
	}

	cfg.maxCharsPerLine = 100
	cfg.maxBufferLines = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero maxBufferLines")
	}

	cfg.maxBufferLines = 1000
	cfg.command = "echo ok"
	cfg.args = []string{"file.log"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for command + args")
	}

	cfg.command = ""
	cfg.args = []string{"a.log", "b.log"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for multiple file args")
	}
}

func TestCfgGetOutputTemplate(t *testing.T) {
	cfg := Cfg{
		outputTemplate: sinkConsoleTemplate.NewTemplateType(),
	}

	if cfg.GetOutputTemplate() != "none" {
		t.Fatalf("unexpected template: %s", cfg.GetOutputTemplate())
	}
}

func TestCfgGetIndicator(t *testing.T) {
	cfg := Cfg{
		indicator: sinkConsoleIndicator.NewIndicatorType(),
	}

	if cfg.GetIndicator() == "" {
		t.Fatal("expected non-empty indicator")
	}
}
