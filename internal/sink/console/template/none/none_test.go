package none

import "testing"

type ind struct{}

func (ind) Clean(bool) string { return "" }
func (ind) Get() string       { return "" }

type win struct{}

func (win) GetPosition() (uint64, uint64)   { return 0, 0 }
func (win) GetBufferSize() (uint64, uint64) { return 0, 0 }
func (win) GetIcon() string                 { return "" }
func (win) GetTitle() string                { return "" }

func TestNoneTemplateMethods(t *testing.T) {
	tpl := New(ind{}, win{})
	if tpl.GetHeader() != "" || tpl.GetCellar() != "" {
		t.Fatal("none template should not render header/cellar")
	}
	if tpl.GetHeaderClean(false) != "" || tpl.GetCellarClean(false) != "" {
		t.Fatal("none template clean methods should be empty")
	}
	if tpl.FormatLine("x") != "x" {
		t.Fatal("FormatLine must keep line")
	}
	if tpl.StartLine() == "" {
		t.Fatal("StartLine should not be empty")
	}
	if tpl.CleanLine() != "" || tpl.EndLine() != "" {
		t.Fatal("none template clean/end should be empty")
	}
}
