package parser_test

import (
	"testing"

	"tail/internal/formatter/buffer/ansi"
	parser "tail/internal/formatter/buffer/parser"
)

type branchCmd struct {
	t ansi.TypeCommand
	u []uint64
	s string
}

func (c branchCmd) GetType() ansi.TypeCommand { return c.t }
func (c branchCmd) GetValsUint() []uint64     { return c.u }
func (c branchCmd) GetValsString() string     { return c.s }

type decBranch struct{}

func (decBranch) GetFunction(ch rune) (func([]uint64) (ansi.Command, bool), bool) {
	if ch == 'h' || ch == 'l' || ch == 'c' {
		return func(v []uint64) (ansi.Command, bool) { return branchCmd{t: ansi.KeyOtherDo, u: v}, true }, true
	}
	return nil, false
}

type secBranch struct{}

func (secBranch) GetFunction(ch rune) (func([]uint64) (ansi.Command, bool), bool) {
	if ch == 'c' || ch == 'm' {
		return func(v []uint64) (ansi.Command, bool) { return branchCmd{t: ansi.KeyTerminalID, u: v}, true }, true
	}
	return nil, false
}

type winBranch struct{}

func (winBranch) GetFunction(ch rune) (func([]uint64, string) (ansi.Command, bool), bool) {
	switch ch {
	case '?':
		return func(v []uint64, _ string) (ansi.Command, bool) { return branchCmd{t: ansi.KeyGetPalettes, u: v}, true }, true
	case 's':
		return func(v []uint64, _ string) (ansi.Command, bool) { return branchCmd{t: ansi.KeySetPalettes, u: v}, true }, true
	case 't':
		return func(v []uint64, s string) (ansi.Command, bool) { return branchCmd{t: ansi.KeyTitle, u: v, s: s}, true }, true
	case 5:
		return func(_ []uint64, s string) (ansi.Command, bool) { return branchCmd{t: ansi.KeyText, s: s}, true }, true
	}
	return nil, false
}

type mainBranch struct{}

func (mainBranch) GetFunction(ch rune) (func([]uint64) (ansi.Command, bool), bool) {
	switch ch {
	case 'E':
		return func(v []uint64) (ansi.Command, bool) { return branchCmd{t: ansi.KeyCursorDownLeft, u: v}, true }, true
	case 'G':
		return func(v []uint64) (ansi.Command, bool) { return branchCmd{t: ansi.KeyCursorSetX, u: v}, true }, true
	case 's':
		return func(v []uint64) (ansi.Command, bool) { return branchCmd{t: ansi.KeyCursorSave, u: v}, true }, true
	case 'u':
		return func(v []uint64) (ansi.Command, bool) { return branchCmd{t: ansi.KeyCursorLoad, u: v}, true }, true
	case 'm':
		return func(v []uint64) (ansi.Command, bool) { return branchCmd{t: ansi.KeyColor, u: v}, true }, true
	}
	return nil, false
}

type otherBranch struct{}

func (otherBranch) GetFunction(rune) (func([]uint64) (ansi.Command, bool), bool) { return nil, false }

func TestParseBranchCoverage(t *testing.T) {
	par := parser.NewParser(decBranch{}, secBranch{}, winBranch{}, mainBranch{}, otherBranch{})

	cases := []string{
		"\x1b[?1h",            // DEC private mode
		"\x1b[>1c",            // Secondary DA
		"\x1b[31m",            // MainList path
		"\x1b]4;1;?X",         // ] ? branch
		"\x1b]4;1;rgb:1;2;3X", // ] rgb branch
		"\x1b]0;1;title\x07X", // ] title branch
		"\x1b7",               // ESC 7 branch
		"\x1b8",               // ESC 8 branch
	}

	for _, in := range cases {
		out := par.Parse(in)
		if len(out) == 0 {
			t.Fatalf("expected commands for input %q", in)
		}
	}
}
