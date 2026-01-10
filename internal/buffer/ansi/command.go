package ansi

type TypeCommand rune

const (
	KeySetCursorXY       TypeCommand = 'H'
	KeyCursorUp          TypeCommand = 'A'
	KeyCursorDown        TypeCommand = 'B'
	KeyCursorRight       TypeCommand = 'C'
	KeyCursorLeft        TypeCommand = 'D'
	KeyCursorDownLeft    TypeCommand = 'E'
	KeyCursorUpLeft      TypeCommand = 'F'
	KeyCursorSetX        TypeCommand = 'G'
	KeyCursorSetY        TypeCommand = 'd'
	KeyCursorSave        TypeCommand = 's'
	KeyCursorLoad        TypeCommand = 'u'
	KeyClearScreen       TypeCommand = 'J'
	KeyClearLine         TypeCommand = 'K'
	KeyColor             TypeCommand = 'm'
	KeyOtherDo           TypeCommand = 'h'
	KeyOtherRevert       TypeCommand = 'l'
	KeyVT100             TypeCommand = 'c'
	KeyGetStatusOrCursor TypeCommand = 'n'
	KeyTerminalID        TypeCommand = 0
	KeyKeysModificator   TypeCommand = 1
	KeyGetPalettes       TypeCommand = 2
	KeySetPalettes       TypeCommand = 3
	KeyTitle             TypeCommand = 4
	KeyText              TypeCommand = 5
)

type Command interface {
	GetType() TypeCommand
	GetValsUint() []uint64
	GetValsString() string
}

type commandUint struct {
	key  TypeCommand
	vals []uint64
}

func (c *commandUint) GetType() TypeCommand {
	return c.key
}
func (c *commandUint) GetValsUint() []uint64 {
	return c.vals
}
func (c *commandUint) GetValsString() string {
	return ""
}

type commandUintString struct {
	key     TypeCommand
	vals    []uint64
	valsStr string
}

func (c *commandUintString) GetType() TypeCommand {
	return c.key
}
func (c *commandUintString) GetValsUint() []uint64 {
	return c.vals
}
func (c *commandUintString) GetValsString() string {
	return c.valsStr
}
