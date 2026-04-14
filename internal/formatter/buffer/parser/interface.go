package parser

import "tail/internal/formatter/buffer/ansi"

type DECPrivateModes interface {
	GetFunction(char rune) (func([]uint64) (ansi.Command, bool), bool)
}
type SecondaryDA interface {
	GetFunction(char rune) (func([]uint64) (ansi.Command, bool), bool)
}
type WindowAndPalettes interface {
	GetFunction(char rune) (func([]uint64, string) (ansi.Command, bool), bool)
}
type MainList interface {
	GetFunction(char rune) (func([]uint64) (ansi.Command, bool), bool)
}
type OtherList interface {
	GetFunction(char rune) (func([]uint64) (ansi.Command, bool), bool)
}

type Command interface {
	GetType() ansi.TypeCommand
	GetValsUint() []uint64
	GetValsString() string
}

type Parser interface {
	Parse(str string) []Command
}
