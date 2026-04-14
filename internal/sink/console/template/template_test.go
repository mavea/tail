package template

import (
	"testing"
)

type testCfg struct {
	outputTemplate string
	indicator      string
}

func (c testCfg) GetOutputTemplate() string { return c.outputTemplate }
func (c testCfg) GetIndicator() string      { return c.indicator }

type testIndicator struct{}

func (testIndicator) Get() string         { return "*" }
func (testIndicator) Clean(_ bool) string { return "" }

type testWindow struct{}

func (testWindow) SetPosition(int, uint64)         {}
func (testWindow) GetPosition() (uint64, uint64)   { return 1, 2 }
func (testWindow) SetBufferSize(uint64, int)       {}
func (testWindow) GetBufferSize() (uint64, uint64) { return 3, 4 }
func (testWindow) SetIcon(string)                  {}
func (testWindow) GetIcon() string                 { return "I" }
func (testWindow) SetTitle(string)                 {}
func (testWindow) GetTitle() string                { return "T" }

func TestTemplateType(t *testing.T) {
	tplType := NewTemplateType()
	if tplType.String() != "none" {
		t.Fatalf("unexpected default template: %s", tplType.String())
	}
	if tplType.Type() != "string" {
		t.Fatalf("unexpected type: %s", tplType.Type())
	}

	if err := tplType.Set("minimal"); err != nil {
		t.Fatalf("unexpected set error: %v", err)
	}
	if !tplType.Validate() {
		t.Fatalf("expected minimal to be valid")
	}

	if err := tplType.Set(""); err != nil || tplType.String() != "none" {
		t.Fatalf("expected empty value to map to none, got %q err=%v", tplType.String(), err)
	}

	if err := tplType.Set("broken"); err == nil {
		t.Fatalf("expected validation error for broken")
	}
}

func TestNewTemplate(t *testing.T) {
	ind := testIndicator{}
	win := testWindow{}

	fullTpl, err := NewTemplate(testCfg{outputTemplate: "full", indicator: "none"}, ind, win)
	if err != nil {
		t.Fatalf("unexpected full template error: %v", err)
	}
	if fullTpl.GetHeader() == "" || fullTpl.GetCellar() == "" {
		t.Fatalf("expected full template to include header and cellar")
	}
	if fullTpl.GetHeaderClean(true) != "" {
		t.Fatalf("expected empty clean header for first line")
	}

	minimalTpl, err := NewTemplate(testCfg{outputTemplate: "minimal", indicator: "none"}, ind, win)
	if err != nil {
		t.Fatalf("unexpected minimal template error: %v", err)
	}
	if minimalTpl.GetHeader() == "" || minimalTpl.GetCellar() != "" {
		t.Fatalf("unexpected minimal template rendering")
	}

	noneTpl, err := NewTemplate(testCfg{outputTemplate: "unknown", indicator: "none"}, ind, win)
	if err != nil {
		t.Fatalf("unexpected none template error: %v", err)
	}
	if noneTpl.GetHeader() != "" || noneTpl.GetCellar() != "" {
		t.Fatalf("none template should not render header/cellar")
	}
}
