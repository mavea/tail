package console

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	formatterWindow "tail/internal/formatter/window"
	"tail/tests/mocks"
)

type failingWriter struct{}

func (f failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

type cfgMock struct{}

func (cfgMock) GetMaxLineCount() int      { return 10 }
func (cfgMock) GetMaxCharsPerLine() int   { return 100 }
func (cfgMock) GetOutputTemplate() string { return "none" }
func (cfgMock) GetIndicator() string      { return "none" }

func TestConsoleNewAndGetDefaultStyle(t *testing.T) {
	var out bytes.Buffer
	writer := bufio.NewWriter(&out)
	c, err := New(context.TODO(), cfgMock{}, writer, formatterWindow.NewWindow("", ""))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if c == nil {
		t.Fatal("expected console instance")
	}
	start, clean, end := c.GetDefaultStyle()
	if start == "" && clean == "" && end == "" {
		// for none template clean/end can be empty, but start must exist
		t.Fatal("unexpected empty default style")
	}
}

func newTemplateMock() *mocks.TemplateMock {
	return &mocks.TemplateMock{
		GetHeaderFunc:      func() string { return "HDR\n" },
		GetHeaderCleanFunc: func(bool) string { return "" },
		GetCellarFunc:      func() string { return "FTR\n" },
		GetCellarCleanFunc: func(bool) string { return "" },
		FormatLineFunc:     func(line string) string { return "[" + line + "]" },
		StartLineFunc:      func() string { return "" },
		EndLineFunc:        func() string { return "" },
		CleanLineFunc:      func() string { return "" },
	}
}

func TestConsolePrintAndSetDataWithGeneratedMock(t *testing.T) {
	var out bytes.Buffer
	c := &console{
		writer:    bufio.NewWriter(&out),
		template:  newTemplateMock(),
		firstLine: true,
	}

	c.SetData([]string{"a", "b"})
	if err := c.Print(); err != nil {
		t.Fatalf("Print error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "[a]") || !strings.Contains(result, "[b]") {
		t.Fatalf("unexpected output: %q", result)
	}
	if !strings.Contains(result, "HDR") || !strings.Contains(result, "FTR") {
		t.Fatalf("header/cellar not rendered: %q", result)
	}
}

func TestConsoleClearScreenWithGeneratedMock(t *testing.T) {
	var out bytes.Buffer
	c := &console{
		writer:    bufio.NewWriter(&out),
		template:  newTemplateMock(),
		firstLine: false,
	}

	if err := c.ClearScreen(); err != nil {
		t.Fatalf("ClearScreen error: %v", err)
	}
}

func TestConsoleErrorAndStatus(t *testing.T) {
	var out bytes.Buffer
	c := &console{
		writer:   bufio.NewWriter(&out),
		template: newTemplateMock(),
	}

	if err := c.Error([]string{"line"}, []string{"err"}); err != nil {
		t.Fatalf("Error print failed: %v", err)
	}
	if err := c.SetStatus(7); err != nil {
		t.Fatalf("SetStatus failed: %v", err)
	}

	r := out.String()
	if !strings.Contains(r, "========= ERROR =========") {
		t.Fatalf("expected error section in output: %q", r)
	}
	if !strings.Contains(r, "========= STATUS =========") {
		t.Fatalf("expected status section in output: %q", r)
	}
}

func TestConsoleErrorWriteFailure(t *testing.T) {
	c := &console{
		writer:   bufio.NewWriter(io.Writer(failingWriter{})),
		template: newTemplateMock(),
	}

	if err := c.Error([]string{"line"}, []string{"err"}); err == nil {
		t.Fatal("expected error on write failure")
	}
}

func TestConsoleSetStatusWriteFailure(t *testing.T) {
	c := &console{
		writer:   bufio.NewWriter(io.Writer(failingWriter{})),
		template: newTemplateMock(),
	}

	if err := c.SetStatus(1); err == nil {
		t.Fatal("expected error on write failure")
	}
}
