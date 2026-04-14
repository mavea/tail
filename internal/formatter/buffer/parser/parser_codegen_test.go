package parser_test

import (
	"testing"

	"tail/internal/formatter/buffer/ansi"
	parser "tail/internal/formatter/buffer/parser"
	"tail/tests/mocks"
)

type cmd struct {
	t ansi.TypeCommand
	u []uint64
	s string
}

func (c cmd) GetType() ansi.TypeCommand { return c.t }
func (c cmd) GetValsUint() []uint64     { return c.u }
func (c cmd) GetValsString() string     { return c.s }

func TestNewParserAndParseUsingGeneratedMocks(t *testing.T) {
	mainMock := &mocks.MainListMock{
		GetFunctionFunc: func(ch rune) (func([]uint64) (ansi.Command, bool), bool) {
			switch ch {
			case 1:
				return func(v []uint64) (ansi.Command, bool) { return cmd{t: ansi.KeyCursorDownLeft, u: v}, true }, true
			case 'E':
				return func(v []uint64) (ansi.Command, bool) { return cmd{t: ansi.KeyCursorDownLeft, u: v}, true }, true
			case 'G':
				return func(v []uint64) (ansi.Command, bool) { return cmd{t: ansi.KeyCursorSetX, u: v}, true }, true
			default:
				return nil, false
			}
		},
	}

	winMock := &mocks.WindowAndPalettesMock{
		GetFunctionFunc: func(ch rune) (func([]uint64, string) (ansi.Command, bool), bool) {
			if ch != 5 {
				return nil, false
			}
			return func(_ []uint64, s string) (ansi.Command, bool) { return cmd{t: ansi.KeyText, s: s}, true }, true
		},
	}

	decMock := &mocks.DECPrivateModesMock{GetFunctionFunc: func(rune) (func([]uint64) (ansi.Command, bool), bool) {
		return nil, false
	}}
	secMock := &mocks.SecondaryDAMock{GetFunctionFunc: func(rune) (func([]uint64) (ansi.Command, bool), bool) {
		return nil, false
	}}
	otherMock := &mocks.OtherListMock{GetFunctionFunc: func(rune) (func([]uint64) (ansi.Command, bool), bool) {
		return nil, false
	}}

	par := parser.NewParser(decMock, secMock, winMock, mainMock, otherMock)
	commands := par.Parse("hello")
	if len(commands) < 2 {
		t.Fatalf("expected text and newline commands, got %d", len(commands))
	}
	if commands[0].GetType() != ansi.KeyText {
		t.Fatalf("expected first command to be text, got %v", commands[0].GetType())
	}
	if commands[0].GetValsString() != "hello" {
		t.Fatalf("unexpected text payload: %q", commands[0].GetValsString())
	}
	if commands[len(commands)-1].GetType() != ansi.KeyCursorDownLeft {
		t.Fatalf("expected last command to be newline, got %v", commands[len(commands)-1].GetType())
	}

	if len(mainMock.GetFunctionCalls()) == 0 {
		t.Fatal("expected main list mock to be called")
	}
	if len(winMock.GetFunctionCalls()) == 0 {
		t.Fatal("expected window/palettes mock to be called")
	}
}
