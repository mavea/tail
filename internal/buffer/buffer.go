package buffer

import (
	"math"
	"sync"

	"tail/internal/buffer/ansi"
	"tail/internal/buffer/line"
	"tail/internal/buffer/parser"
	"tail/internal/buffer/style"
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
}

func New(cfg config) Buffer {
	return &buffer{
		cfg:    cfg,
		text:   line.NewLines(1),
		parser: parser.NewParser(ansi.GetLists()),
		style:  style.New(),
	}
}

func (buf *buffer) execCommands(comm ansi.Command) {
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
		if val[0] == 2 {
			buf.cfg.SetProcessName(comm.GetValsString())
		}
	case ansi.KeyTerminalID:
		buf.x = buf.text.Get(buf.y).Set(buf.style.String(), buf.cfg.GetProcessName(), buf.x)

	case ansi.KeyText:
		buf.x = buf.text.Get(buf.y).Set(buf.style.String(), comm.GetValsString(), buf.x)
	case ansi.KeySetCursorXY:
		val := comm.GetValsUint()
		buf.y = val[0]

		if val[0] > math.MaxInt {
			buf.x = math.MaxInt
		} else {
			// #nosec G115
			buf.x = int(val[1])
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
	case ansi.KeyCursorDown, ansi.KeyCursorDownLeft:
		val := comm.GetValsUint()
		buf.y += val[0]
		if buf.y >= buf.text.Len() {
			buf.text.Add(buf.y - buf.text.Len() + 1)
		}
		if comm.GetType() == ansi.KeyCursorDownLeft {
			buf.x = 0
		}
	case ansi.KeyCursorRight:
		val := comm.GetValsUint()
		if val[0] > math.MaxInt || val[0]+uint64(buf.x) > math.MaxInt {
			buf.x = math.MaxInt
		} else {
			// #nosec G115
			buf.x += int(val[0])
		}
	case ansi.KeyCursorLeft:
		val := comm.GetValsUint()
		if val[0] > math.MaxInt {
			buf.x = 0
		} else {
			// #nosec G115
			buf.x -= int(val[0])
			if buf.x < 0 {
				buf.x = 0
			}
		}
	case ansi.KeyCursorSetX:
		val := comm.GetValsUint()
		if val[0] > math.MaxInt {
			buf.x = math.MaxInt
		} else {
			// #nosec G115
			buf.x = int(val[0])
		}
	case ansi.KeyCursorSetY:
		val := comm.GetValsUint()
		buf.y = val[0]
		if buf.y >= buf.text.Len() {
			buf.text.Add(buf.y - buf.text.Len() + 1)
		}
	case ansi.KeyCursorSave:
		buf.saveX = buf.x
		buf.saveY = buf.y
	case ansi.KeyCursorLoad:
		buf.x = buf.saveX
		buf.y = buf.saveY
		if buf.y >= buf.text.Len() {
			buf.text.Add(buf.y - buf.text.Len() + 1)
		}
	case ansi.KeyClearScreen:
		val := comm.GetValsUint()
		switch val[0] {
		case 0:
			// e[0J — очистить от курсора до конца экрана
			buf.x = buf.text.Get(buf.y).CleanPostfix(buf.x)
			buf.text.CleanPostfix(buf.y + 1)
		case 1:
			// e[1J — очистить от начала экрана до курсора
			buf.x = buf.text.Get(buf.y).CleanPrefix(buf.x)
			buf.text.CleanPrefix(buf.y)
			buf.y = 0
			buf.x = 0
		case 2:
			// e[2J — очистить весь экран (курсор в 0,0)
			buf.x = 0
			buf.y = 0
			buf.text = line.NewLines(1)
		case 3:
			// e[3J — очистить весь экран и буфер прокрутки
			buf.x = 0
			buf.y = 0
			buf.saveX = 0
			buf.saveY = 0
			buf.text = line.NewLines(1)
			buf.style = style.New()
		}
	case ansi.KeyClearLine:
		val := comm.GetValsUint()
		switch val[0] {
		case 0:
			// e[0K — очистить от курсора до конца строки
			buf.x = buf.text.Get(buf.y).CleanPostfix(buf.x)
		case 1:
			// e[1K — очистить от начала строки до курсора
			buf.x = buf.text.Get(buf.y).CleanPrefix(buf.x)
			buf.x = 0
		case 2:
			// e[2K — очистить всю строку
			buf.text.CleanString(buf.y)
			buf.x = 0
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

func (buf *buffer) GetLast(len uint64) []string {
	buf.mu.RLock()
	defer buf.mu.RUnlock()

	switch true {
	case buf.text.Len() == 0:
		return []string{}
	case buf.text.Len() < len:
		arr := make([]string, 0, buf.text.Len())
		for i := uint64(0); i < buf.text.Len(); i++ {
			arr = append(arr, buf.startLine+buf.text.Get(i).String(buf.cleanLine, buf.cfg.GetLengthLines())+buf.endLine)
		}

		return arr
	default:
		arr := make([]string, 0, len)
		for i := uint64(buf.text.Len() - len); i < buf.text.Len(); i++ {
			arr = append(arr, buf.startLine+buf.text.Get(i).String(buf.cleanLine, buf.cfg.GetLengthLines())+buf.endLine)
		}

		return arr
	}
}

func (buf *buffer) GetFull() []string {
	buf.mu.RLock()
	defer buf.mu.RUnlock()
	arr := make([]string, 0, buf.text.Len())
	for i := uint64(0); i < buf.text.Len(); i++ {
		arr = append(arr, buf.startLine+buf.text.Get(i).String(buf.cleanLine, buf.cfg.GetLengthLines())+buf.endLine)
	}

	return arr
}
