package bootstrap

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	langEn "tail/internal/lang/en"
)

type testBootstrapCfg struct {
	help    bool
	version bool
}

func (c testBootstrapCfg) GetMaxLineCount() int      { return 10 }
func (c testBootstrapCfg) GetMaxCharsPerLine() int   { return 100 }
func (c testBootstrapCfg) GetMaxBufferLines() uint64 { return 1000 }
func (c testBootstrapCfg) GetProcessName() string    { return "" }
func (c testBootstrapCfg) GetProcessIcon() string    { return "" }
func (c testBootstrapCfg) GetOutputTemplate() string { return "none" }
func (c testBootstrapCfg) GetIndicator() string      { return "none" }
func (c testBootstrapCfg) GetOutputMode() string     { return "direct" }
func (c testBootstrapCfg) IsHelp() bool              { return c.help }
func (c testBootstrapCfg) IsVersion() bool           { return c.version }
func (c testBootstrapCfg) GetCommand() string        { return "" }
func (c testBootstrapCfg) GetArgs() []string         { return nil }
func (c testBootstrapCfg) IsCSIEnabled() bool        { return true }
func (c testBootstrapCfg) IsFullOutput() bool        { return false }

func TestNewRunModeHelp(t *testing.T) {
	cfg := testBootstrapCfg{help: true}
	lang := langEn.NewLang()
	stdin := os.Stdin
	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)

	rm := NewRunMode(cfg, lang, stdin, stdout, stderr)
	if rm == nil {
		t.Fatal("expected RunMode instance")
	}

	route := rm.Route()
	if route == nil {
		t.Fatal("expected route function")
	}
}

func TestNewRunModeVersion(t *testing.T) {
	cfg := testBootstrapCfg{version: true}
	lang := langEn.NewLang()
	stdin := os.Stdin
	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)

	rm := NewRunMode(cfg, lang, stdin, stdout, stderr)
	route := rm.Route()

	if route == nil {
		t.Fatal("expected route function")
	}
}

func TestNewRunModeTail(t *testing.T) {
	cfg := testBootstrapCfg{}
	lang := langEn.NewLang()
	stdin := os.Stdin
	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)

	rm := NewRunMode(cfg, lang, stdin, stdout, stderr)
	route := rm.Route()

	if route == nil {
		t.Fatal("expected route function")
	}
}

func TestHelpOutput(t *testing.T) {
	cfg := testBootstrapCfg{help: true}
	lang := langEn.NewLang()
	stdin := os.Stdin
	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)

	rm := NewRunMode(cfg, lang, stdin, stdout, stderr)
	ctx := context.Background()

	err := rm.help(ctx)
	if err != nil {
		t.Logf("help returned error: %v", err)
	}
}

func TestVersionOutput(t *testing.T) {
	cfg := testBootstrapCfg{version: true}
	lang := langEn.NewLang()
	stdin := os.Stdin
	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)

	rm := NewRunMode(cfg, lang, stdin, stdout, stderr)
	ctx := context.Background()

	err := rm.version(ctx)
	if err != nil {
		t.Logf("version returned error: %v", err)
	}
}

func TestVersionOutputUsesBuildVariables(t *testing.T) {
	origVersion := Version
	origBuildTime := BuildTime
	defer func() {
		Version = origVersion
		BuildTime = origBuildTime
	}()

	Version = "1.2.3"
	BuildTime = "2026-04-12"

	var out bytes.Buffer
	rm := NewRunMode(
		testBootstrapCfg{version: true},
		langEn.NewLang(),
		os.Stdin,
		bufio.NewWriter(&out),
		bufio.NewWriter(os.Stderr),
	)

	if err := rm.version(context.Background()); err != nil {
		t.Fatalf("version returned error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "1.2.3") {
		t.Fatalf("expected version in output, got: %q", got)
	}
	if !strings.Contains(got, "2026-04-12") {
		t.Fatalf("expected build time in output, got: %q", got)
	}
}

func TestHelpOutputNilStdout(t *testing.T) {
	rm := NewRunMode(testBootstrapCfg{help: true}, langEn.NewLang(), os.Stdin, nil, bufio.NewWriter(os.Stderr))
	if err := rm.help(context.Background()); err == nil {
		t.Fatal("expected help error for nil stdout")
	}
}

func TestVersionOutputNilStdout(t *testing.T) {
	rm := NewRunMode(testBootstrapCfg{version: true}, langEn.NewLang(), os.Stdin, nil, bufio.NewWriter(os.Stderr))
	if err := rm.version(context.Background()); err == nil {
		t.Fatal("expected version error for nil stdout")
	}
}

func TestBootstrapCfgInterface(t *testing.T) {
	cfg := testBootstrapCfg{}
	_ = Cfg(cfg)
}
