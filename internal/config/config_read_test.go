package config

import (
	"os"
	"tail/internal/formatter"
	"testing"

	"github.com/spf13/cobra"

	langEn "tail/internal/lang/en"
	sinkConsoleIndicator "tail/internal/sink/console/indicator"
	sinkConsoleTemplate "tail/internal/sink/console/template"
)

func TestReadConfAndReadArguments(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"tail", "-n", "10", "-l", "50", "-s", "1000", "-a", "proc", "-i", "*", "-o", "thread", "-t", "minimal", "-r", "roller", "-c", "echo hi", "file.txt"}

	cfg, err := ReadConf(langEn.NewLang())
	if err != nil {
		t.Fatalf("ReadConf error: %v", err)
	}

	if cfg.GetMaxLineCount() != 10 || cfg.GetMaxCharsPerLine() != 50 {
		t.Fatalf("unexpected numeric flags: n=%d l=%d", cfg.GetMaxLineCount(), cfg.GetMaxCharsPerLine())
	}
	if cfg.GetOutputMode() != "thread" {
		t.Fatalf("unexpected output mode: %s", cfg.GetOutputMode())
	}
	if cfg.GetOutputTemplate() != "minimal" {
		t.Fatalf("unexpected output template: %s", cfg.GetOutputTemplate())
	}
	if cfg.GetIndicator() != "roller" {
		t.Fatalf("unexpected indicator: %s", cfg.GetIndicator())
	}
	if cfg.GetCommand() != "echo hi" {
		t.Fatalf("unexpected command: %s", cfg.GetCommand())
	}
	if len(cfg.GetArgs()) != 1 || cfg.GetArgs()[0] != "file.txt" {
		t.Fatalf("unexpected args: %+v", cfg.GetArgs())
	}
}

func TestReadArgumentsExecuteError(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"tail"}

	cfg := &Cfg{
		outputMode:     formatter.NewOutputMode(),
		outputTemplate: sinkConsoleTemplate.NewTemplateType(),
		indicator:      sinkConsoleIndicator.NewIndicatorType(),
	}
	cmd := &cobra.Command{
		Use: "tail",
		RunE: func(cmd *cobra.Command, args []string) error {
			return os.ErrInvalid
		},
	}

	_, err := cfg.ReadArguments(cmd, cfg, langEn.NewLang())
	if err == nil {
		t.Fatal("expected execute error")
	}
}
