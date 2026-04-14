package full

import "testing"

type ind struct{}

func (ind) Clean(bool) string { return "C" }
func (ind) Get() string       { return "I" }

type win struct{}

func (win) GetPosition() (uint64, uint64)   { return 1, 2 }
func (win) GetBufferSize() (uint64, uint64) { return 3, 4 }
func (win) GetIcon() string                 { return "*" }
func (win) GetTitle() string                { return "title" }

func TestFullTemplateMethods(t *testing.T) {
	tpl := New(ind{}, win{})
	if tpl.GetHeader() == "" {
		t.Fatal("header should not be empty")
	}
	if tpl.GetHeaderClean(false) == "" {
		t.Fatal("header clean should not be empty for non-first line")
	}
	if tpl.GetCellar() == "" {
		t.Fatal("cellar should not be empty")
	}
	if tpl.GetHeaderClean(true) != "" {
		t.Fatal("header clean for first line should be empty")
	}
	if tpl.GetCellarClean(false) == "" {
		t.Fatal("cellar clean should not be empty for non-first line")
	}
	if tpl.GetCellarClean(true) != "" {
		t.Fatal("cellar clean for first line should be empty")
	}
	if tpl.FormatLine("x") != "x" {
		t.Fatal("FormatLine must keep line")
	}
	if tpl.StartLine() == "" {
		t.Fatal("StartLine should not be empty")
	}
	if tpl.CleanLine() != "" {
		t.Fatal("CleanLine should be empty")
	}
	if tpl.EndLine() != "" {
		t.Fatal("EndLine should be empty")
	}
}
