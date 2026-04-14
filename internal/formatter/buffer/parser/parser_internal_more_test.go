package parser

import (
	"testing"

	"tail/internal/formatter/buffer/ansi"
)

type localCommand struct {
	t ansi.TypeCommand
	u []uint64
	s string
}

func (c localCommand) GetType() ansi.TypeCommand { return c.t }
func (c localCommand) GetValsUint() []uint64     { return c.u }
func (c localCommand) GetValsString() string     { return c.s }

type localDEC struct{}

func (localDEC) GetFunction(ch rune) (func([]uint64) (ansi.Command, bool), bool) {
	if ch == 'h' || ch == 'l' || ch == 'c' {
		return func(v []uint64) (ansi.Command, bool) { return localCommand{t: ansi.KeyOtherDo, u: v}, true }, true
	}
	return nil, false
}

type localSecondary struct{}

func (localSecondary) GetFunction(ch rune) (func([]uint64) (ansi.Command, bool), bool) {
	if ch == 'c' || ch == 'm' {
		return func(v []uint64) (ansi.Command, bool) { return localCommand{t: ansi.KeyTerminalID, u: v}, true }, true
	}
	return nil, false
}

type localWindow struct{}

func (localWindow) GetFunction(ch rune) (func([]uint64, string) (ansi.Command, bool), bool) {
	switch ch {
	case '?':
		return func(v []uint64, _ string) (ansi.Command, bool) {
			return localCommand{t: ansi.KeyGetPalettes, u: v}, true
		}, true
	case 's':
		return func(v []uint64, _ string) (ansi.Command, bool) {
			return localCommand{t: ansi.KeySetPalettes, u: v}, true
		}, true
	case 't':
		return func(v []uint64, s string) (ansi.Command, bool) {
			return localCommand{t: ansi.KeyTitle, u: v, s: s}, true
		}, true
	case 5:
		return func(_ []uint64, s string) (ansi.Command, bool) { return localCommand{t: ansi.KeyText, s: s}, true }, true
	}
	return nil, false
}

type localMain struct{}

func (localMain) GetFunction(ch rune) (func([]uint64) (ansi.Command, bool), bool) {
	switch ch {
	case 'E':
		return func(v []uint64) (ansi.Command, bool) { return localCommand{t: ansi.KeyCursorDownLeft, u: v}, true }, true
	case 'G':
		return func(v []uint64) (ansi.Command, bool) { return localCommand{t: ansi.KeyCursorSetX, u: v}, true }, true
	case 's':
		return func(v []uint64) (ansi.Command, bool) { return localCommand{t: ansi.KeyCursorSave, u: v}, true }, true
	case 'u':
		return func(v []uint64) (ansi.Command, bool) { return localCommand{t: ansi.KeyCursorLoad, u: v}, true }, true
	case 'm':
		return func(v []uint64) (ansi.Command, bool) { return localCommand{t: ansi.KeyColor, u: v}, true }, true
	}
	return nil, false
}

type localMainEmpty struct{}

func (localMainEmpty) GetFunction(rune) (func([]uint64) (ansi.Command, bool), bool) {
	return nil, false
}

type localOther struct{}

func (localOther) GetFunction(rune) (func([]uint64) (ansi.Command, bool), bool) { return nil, false }

func TestEnterInParseANSIEscapeSequencesBranches(t *testing.T) {
	par := &parser{
		ansiDECPrivateModes:   localDEC{},
		ansiSecondaryDA:       localSecondary{},
		ansiWindowAndPalettes: localWindow{},
		ansiMainList:          localMain{},
		ansiOtherList:         localOther{},
	}

	cases := []string{
		"\x1b[?1h",
		"\x1b[>1c",
		"\x1b[31m",
		"\x1b]4;1;?X",
		"\x1b]4;1;rgb:1;2;3X",
		"\x1b]0;1;title\x07X",
		"\x1b7",
		"\x1b8",
	}
	for _, tc := range cases {
		i, cmd := par.enterInParseANSIEscapeSequences([]rune(tc))
		if i < 1 {
			t.Fatalf("expected positive offset for %q, got %d", tc, i)
		}
		if cmd == nil {
			t.Fatalf("expected command for %q", tc)
		}
	}

	// malformed/unsupported should not panic and should return nil command
	if i, cmd := par.enterInParseANSIEscapeSequences([]rune("\x1b")); i != 1 || cmd != nil {
		t.Fatalf("unexpected result for short escape: i=%d cmd=%v", i, cmd)
	}
	if i, cmd := par.enterInParseANSIEscapeSequences([]rune("\x1bx")); i < 1 || cmd != nil {
		t.Fatalf("unexpected result for unknown escape: i=%d cmd=%v", i, cmd)
	}
}

func TestGetLineCommandsNilWhenMainListHasNoHandlers(t *testing.T) {
	par := &parser{ansiMainList: localMainEmpty{}}
	if cmd := par.getNewLineCommand(); cmd != nil {
		t.Fatal("expected nil new line command")
	}
	if cmd := par.getBeginningOfLineCommand(); cmd != nil {
		t.Fatal("expected nil beginning-of-line command")
	}
}
