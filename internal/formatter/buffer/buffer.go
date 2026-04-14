package buffer

import (
	"math"
	"sync"

	"tail/internal/formatter/buffer/ansi"
	"tail/internal/formatter/buffer/line"
	"tail/internal/formatter/buffer/parser"
	"tail/internal/formatter/buffer/style"
)

type buffer struct {
	cfg    config
	parser parser.Parser
	mu     sync.RWMutex

	text  line.Lines
	style style.Style

	startLine string
	cleanLine string
	endLine   string

	x     int
	saveX int
	y     uint64
	saveY uint64

	window Window
}

func uint64ToIntClamp(v uint64) int {
	if v >= uint64(math.MaxInt) {
		return math.MaxInt
	}

	return int(v)
}

func addUintToIntClamp(base int, add uint64) int {
	if base < 0 {
		base = 0
	}
	if add >= uint64(math.MaxInt-base) {
		return math.MaxInt
	}

	// #nosec G115 -- add is bounded by math.MaxInt-base above.
	return base + int(add)
}

func subUintFromIntFloor(base int, sub uint64) int {
	if base <= 0 {
		return 0
	}
	if sub >= uint64(base) {
		return 0
	}

	// #nosec G115 -- sub is strictly less than base above, so conversion is safe.
	return base - int(sub)
}

func New(cfg config, window Window) Buffer {
	return &buffer{
		cfg:    cfg,
		text:   line.MakeLines(cfg.GetMaxBufferLines()),
		parser: parser.NewParser(ansi.GetLists()),
		style:  style.New(),
		window: window,
	}
}

func (buf *buffer) execCommands(comm ansi.Command) {
	defer buf.window.SetPosition(buf.x, buf.y)
	switch comm.GetType() {
	case ansi.KeyOtherDo:
		val := comm.GetValsUint()
		if val[0] == 1048 || val[0] == 1049 {
			buf.saveX = buf.x
			buf.saveY = buf.y
		}
	case ansi.KeyOtherRevert, ansi.KeyVT100, ansi.KeyKeysModificator:
		// @ не нужные команды
	case ansi.KeyGetPalettes, ansi.KeySetPalettes:
		// @ палитра. Временно не делаем
	case ansi.KeyTitle:
		//Заголовок окна
		//		ESC]0;<текст>\x07 — заголовок и иконки
		//		ESC]1;<текст>\x07 — заголовок иконки
		//		ESC]2;<текст>\x07 — заголовок окна
		val := comm.GetValsUint()
		if len(val) > 0 {
			switch val[0] {
			case 0, 2:
				buf.window.SetTitle(comm.GetValsString())
			case 1:
				buf.window.SetIcon(comm.GetValsString())
			}
		}
	case ansi.KeyTerminalID:
		buf.x = buf.text.Get(buf.y).Set(buf.style.String(), buf.cfg.GetProcessName(), buf.x)
		buf.window.SetBufferSize(buf.text.LenHistory(), buf.x)
	case ansi.KeyText:
		buf.x = buf.text.Get(buf.y).Set(buf.style.String(), comm.GetValsString(), buf.x)
		buf.window.SetBufferSize(buf.text.LenHistory(), buf.x)
	case ansi.KeySetCursorXY:
		val := comm.GetValsUint()
		buf.y = val[0]
		if buf.y > 0 {
			buf.y--
		}
		if buf.y >= buf.text.LenHistory() {
			buf.y = buf.text.LenHistory() - 1
		}
		buf.x = uint64ToIntClamp(val[1])
		if buf.x > 0 {
			buf.x--
		}
	case ansi.KeyCursorUp, ansi.KeyCursorUpLeft:
		val := comm.GetValsUint()
		if buf.y >= val[0] {
			buf.y -= val[0]
		} else {
			buf.y = 0
		}
		if comm.GetType() == ansi.KeyCursorUpLeft {
			buf.x = 0
		}
	case ansi.KeyNewLine:
		buf.y++
		if buf.y >= buf.text.LenHistory() {
			buf.text.Add(1)
			buf.y = buf.text.LenHistory() - 1
		}
		buf.x = 0
	case ansi.KeyCursorDown, ansi.KeyCursorDownLeft:
		val := comm.GetValsUint()
		buf.y += val[0]
		if buf.y >= buf.text.LenHistory() {
			buf.y = buf.text.LenHistory() - 1
		}
		if comm.GetType() == ansi.KeyCursorDownLeft {
			buf.x = 0
		}
		buf.window.SetBufferSize(buf.text.LenHistory(), buf.x)
	case ansi.KeyCursorRight:
		val := comm.GetValsUint()
		buf.x = addUintToIntClamp(buf.x, val[0])
	case ansi.KeyCursorLeft:
		val := comm.GetValsUint()
		buf.x = subUintFromIntFloor(buf.x, val[0])
	case ansi.KeyCursorSetX:
		val := comm.GetValsUint()
		buf.x = uint64ToIntClamp(val[0])
		if buf.x > 0 {
			buf.x--
		}
	case ansi.KeyCursorSetY:
		val := comm.GetValsUint()
		buf.y = val[0]
		if buf.y > 0 {
			buf.y--
		}
		if buf.y >= buf.text.LenHistory() {
			buf.text.Add(buf.y - buf.text.LenHistory() + 1)
		}
		buf.window.SetBufferSize(buf.text.LenHistory(), 0)
	case ansi.KeyCursorSave:
		if !buf.cfg.IsCSIEnabled() {
			return
		}
		buf.saveX = buf.x
		buf.saveY = buf.y
	case ansi.KeyCursorLoad:
		if !buf.cfg.IsCSIEnabled() {
			return
		}
		buf.x = buf.saveX //todo в консоли почему то не работает
		buf.y = buf.saveY
		if buf.y >= buf.text.LenHistory() {
			buf.text.Add(buf.y - buf.text.LenHistory() + 1)
			buf.window.SetBufferSize(buf.text.LenHistory(), 0)
		}
	case ansi.KeyClearScreen:
		val := comm.GetValsUint()
		switch val[0] {
		case 0:
			// e[0J — очистить от курсора до конца экрана
			buf.x = buf.text.Get(buf.y).CleanPostfix(buf.x)
			buf.text.CleanPostfix(buf.y + 1)
			buf.window.SetBufferSize(buf.text.LenHistory(), buf.x)
		case 1:
			// e[1J — очистить от начала экрана до курсора
			buf.x = buf.text.Get(buf.y).CleanPrefix(buf.x)
			buf.text.CleanPrefix(buf.y)
			buf.y = 0
			buf.x = 0
			buf.window.SetBufferSize(buf.text.LenHistory(), 0)
		case 2:
			buf.x = 0
			if !buf.cfg.IsCSIEnabled() {
				buf.text.Add(buf.window.Height())
				buf.y = buf.text.LenHistory() - 1

				return
			}
			// e[2J — очистить весь экран (курсор в 0,0)
			buf.y = 0
			buf.text = line.MakeLines(buf.cfg.GetMaxBufferLines())
			buf.window.SetBufferSize(buf.text.LenHistory(), 0)
		case 3:
			// e[3J — очистить весь экран и буфер прокрутки
			buf.x = 0
			buf.y = 0
			buf.saveX = 0
			buf.saveY = 0
			buf.text = line.MakeLines(buf.cfg.GetMaxBufferLines())
			buf.window.SetBufferSize(buf.text.LenHistory(), 0)
			buf.style = style.New()
		}
	case ansi.KeyClearLine:
		val := comm.GetValsUint()
		switch val[0] {
		case 0:
			// e[0K — очистить от курсора до конца строки
			buf.x = buf.text.Get(buf.y).CleanPostfix(buf.x)
			buf.window.SetBufferSize(buf.text.LenHistory(), buf.x)
		case 1:
			// e[1K — очистить от начала строки до курсора
			buf.x = buf.text.Get(buf.y).CleanPrefix(buf.x)
			buf.x = 0
			buf.window.SetBufferSize(buf.text.LenHistory(), 0)
		case 2:
			// e[2K — очистить всю строку
			buf.text.CleanString(buf.y)
			buf.x = 0
			buf.window.SetBufferSize(buf.text.LenHistory(), 0)
		}
	case ansi.KeyColor:
		val := comm.GetValsUint()
		buf.style = buf.style.Set(val)
	}
}

func (buf *buffer) Add(data string) {
	commands := buf.parser.Parse(data)
	buf.mu.Lock()
	defer buf.mu.Unlock()
	for _, comm := range commands {
		buf.execCommands(comm)
	}
}

func (buf *buffer) SetDefaultStyle(start string, clean string, end string) {
	buf.startLine = start
	buf.cleanLine = clean
	buf.endLine = end
}

func (buf *buffer) GetLast(count int) []string {
	buf.mu.RLock()
	defer buf.mu.RUnlock()
	if count <= 0 || buf.text.LenHistory() == 0 {
		return []string{}
	}

	data := buf.text.GetLastLines(count)
	arr := make([]string, 0, len(data))
	for _, dataLine := range data {
		arr = append(arr, buf.startLine+dataLine.String(buf.cleanLine, buf.cfg.GetMaxCharsPerLine())+buf.endLine)
	}

	return arr
}

func (buf *buffer) GetFull() []string {
	buf.mu.RLock()
	defer buf.mu.RUnlock()
	data := buf.text.GetFullLines()

	arr := make([]string, 0, len(data))
	for _, dataLine := range data {
		arr = append(arr, buf.startLine+dataLine.String(buf.cleanLine, buf.cfg.GetMaxCharsPerLine())+buf.endLine)
	}

	return arr
}
