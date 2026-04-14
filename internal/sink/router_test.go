package sink

import (
	"bufio"
	"bytes"
	"context"
	"testing"

	formatterWindow "tail/internal/formatter/window"
)

type testSinkCfg struct {
	maxLineCount    int
	maxCharsPerLine int
	maxBufferLines  uint64
	outputTemplate  string
	indicator       string
	processName     string
	processIcon     string
	outputMode      string
}

func (c testSinkCfg) GetMaxLineCount() int      { return c.maxLineCount }
func (c testSinkCfg) GetMaxCharsPerLine() int   { return c.maxCharsPerLine }
func (c testSinkCfg) GetOutputTemplate() string { return c.outputTemplate }
func (c testSinkCfg) GetIndicator() string      { return c.indicator }
func (c testSinkCfg) GetMaxBufferLines() uint64 { return c.maxBufferLines }
func (c testSinkCfg) GetProcessName() string    { return c.processName }
func (c testSinkCfg) GetProcessIcon() string    { return c.processIcon }
func (c testSinkCfg) GetOutputMode() string     { return c.outputMode }

func TestNewWriter(t *testing.T) {
	cfg := testSinkCfg{
		maxLineCount:    5,
		maxCharsPerLine: 100,
		outputTemplate:  "none",
		indicator:       "none",
	}
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	window := formatterWindow.NewWindow("", "")

	target, err := NewWriter(context.Background(), cfg, writer, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target == nil {
		t.Fatal("expected target instance")
	}
}

func TestNewWriterWithFull(t *testing.T) {
	cfg := testSinkCfg{
		maxLineCount:    5,
		maxCharsPerLine: 100,
		outputTemplate:  "full",
		indicator:       "roller",
	}
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	window := formatterWindow.NewWindow("🚀", "test")

	target, err := NewWriter(context.Background(), cfg, writer, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target == nil {
		t.Fatal("expected target instance")
	}
}

func TestNewWriterWithMinimal(t *testing.T) {
	cfg := testSinkCfg{
		maxLineCount:    5,
		maxCharsPerLine: 100,
		outputTemplate:  "minimal",
		indicator:       "none",
	}
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	window := formatterWindow.NewWindow("", "")

	target, err := NewWriter(context.Background(), cfg, writer, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target == nil {
		t.Fatal("expected target instance")
	}
}
