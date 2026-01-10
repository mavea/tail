package parser

import "strconv"

type parser struct {
	ansiDECPrivateModes   DECPrivateModes
	ansiSecondaryDA       SecondaryDA
	ansiWindowAndPalettes WindowAndPalettes
	ansiMainList          MainList
	ansiOtherList         OtherList
}

func NewParser(
	ansiDECPrivateModes DECPrivateModes,
	ansiSecondaryDA SecondaryDA,
	ansiWindowAndPalettes WindowAndPalettes,
	ansiMainList MainList,
	ansiOtherList OtherList,
) Parser {
	par := &parser{
		ansiDECPrivateModes:   ansiDECPrivateModes,
		ansiSecondaryDA:       ansiSecondaryDA,
		ansiWindowAndPalettes: ansiWindowAndPalettes,
		ansiMainList:          ansiMainList,
		ansiOtherList:         ansiOtherList,
	}

	return par
}

func (par *parser) parseUintSlice(strRunes []rune) ([]uint64, int) {
	var (
		i     int
		start int
		num   uint64
		rez   = make([]uint64, 0)
		err   error
	)

	for ; i < len(strRunes); i++ {
		num = 0
		start = i
		for ; i < len(strRunes) && strRunes[i] >= '0' && strRunes[i] <= '9'; i++ {
		}
		if i > start {
			if num, err = strconv.ParseUint(string(strRunes[start:i]), 10, 64); err != nil {
				return nil, 0
			}
		}
		rez = append(rez, num)
		if i >= len(strRunes) || strRunes[i] != ';' {
			break
		}
	}

	return rez, i
}

func (par *parser) parseString(strRunes []rune) (string, int) {
	if len(strRunes) == 0 {
		return "", 0
	}

	var (
		i     = 0
		start = i
	)
	for ; i < len(strRunes) && strRunes[i] != '\x07'; i++ {
		if strRunes[i] == '\x1b' {
			return "", 0
		}
	}

	if i < len(strRunes) {
		return string(strRunes[start:i]), i + 1
	}

	return "", 0
}

func (par *parser) getTextCommand(str []rune, start, end int) Command {
	if start < end && len(str) >= end {
		if f, b := par.ansiWindowAndPalettes.GetFunction(5); b {
			if command, bb := f(nil, string(str[start:end])); bb {
				return command
			}
		}
	}

	return nil
}
func (par *parser) getNewLineCommand() Command {
	if f, b := par.ansiMainList.GetFunction('E'); b {
		if command, bb := f([]uint64{1}); bb {
			return command
		}
	}

	return nil
}
func (par *parser) getBeginningOfLineCommand() Command {
	if f, b := par.ansiMainList.GetFunction('G'); b {
		if command, bb := f([]uint64{0}); bb {
			return command
		}
	}

	return nil
}

func (par *parser) enterInParseANSIEscapeSequences(strRunes []rune) (int, Command) {
	var (
		command Command
		i       = 1
		ok      bool
	)
	if len(strRunes) <= i {
		return 1, nil
	}
	switch strRunes[i] {
	case '[': //ESC[?<c>, ESC[<n><c>, ESC[<c>, ESC[><c>, ESC[><n><c>
		i++
		if len(strRunes) <= i {
			return 1, nil
		}
		switch strRunes[i] {
		case '?': //  [h][l][c]
			//Установка режимов (Set Mode)
			//		`ESC[?<n>h`
			//				1 — курсор-ключи (DECCKM)
			//				3 — 80/132 колонок (DECCOLM)
			//				5 — инвертированный экран (DECSCNM)
			//				6 — относительное/абсолютное origin (DECOM)
			//				7 — автоперенос (DECAWM)
			//				8 — авто-повтор клавиш (DECARM)
			//				12 — мигающий курсор (ATT610)
			//				25 — видимый курсор (DECTCEM)
			//				47 — альтернативный буфер
			//				1047 — альтернативный буфер (альтернатива)
			//				1048 — сохранить курсор (альтернатива)
			//				1049 — alt buffer + сохранить курсор
			//Сброс режимов (Reset Mode)
			//		`ESC[?<n>l` (буква L)
			//				(Те же номера, что и выше)
			//Мышь
			//		`ESC[?1000h` — включить мышь (X10)
			//		`ESC[?1002h` — включить мышь (для drag)
			//		`ESC[?1003h` — включить мышь (все события)
			//		`ESC[?1006h` — SGR-формат координат мыши
			//7. Шрифты и клавиатура
			//		`ESC[?1c` — ответ: ESC[?1;0c (базовый VT100)
			//Фокус
			//		`ESC[?1004h` — отслеживание фокуса
			i++
			if len(strRunes) <= i {
				return 1, nil
			}
			sl, offset := par.parseUintSlice(strRunes[i:])
			i += offset
			if len(sl) == 0 || len(strRunes) <= i {
				return 1, nil
			}
			if f, b := par.ansiDECPrivateModes.GetFunction(strRunes[i]); b {
				if command, ok = f(sl); !ok {
					return 1, nil
				}
				return i + 1, command
			} else {
				return 1, nil
			}
		case '>': //[m][c]
			//7. Шрифты и клавиатура
			//		`ESC[>c` — запросить идентификатор терминала
			//		`ESC[>0;256m` — модификаторы клавиш (xterm)
			i++
			if len(strRunes) <= i {
				return 1, nil
			}
			sl, offset := par.parseUintSlice(strRunes[i:])
			i += offset
			if len(sl) == 0 || len(strRunes) <= i {
				return 1, nil
			}
			if f, b := par.ansiSecondaryDA.GetFunction(strRunes[i]); b {
				if command, ok = f(sl); !ok {
					return 1, nil
				}
				return i + 1, command
			} else {
				return 1, nil
			}
		default: //
			//Позиционирование
			//v		`ESC[<n>A` — курсор вверх на n строк (Cursor Up, CUU)
			//		`ESC[<n>B` — курсор вниз на n строк (Cursor Down, CUD)
			//		`ESC[<n>C` — курсор вправо на n символов (Cursor Forward, CUF)
			//		`ESC[<n>D` — курсор влево на n символов (Cursor Backward, CUB)
			//		`ESC[<n>E` — курсор вниз на n строк, в начало строки (CNL)
			//		`ESC[<n>F` — курсор вверх на n строк, в начало строки (CPL)
			//Абсолютное позиционирование
			//		`ESC[<row>;<column>H` — переместить курсор (Cursor Position, CUP)
			//		`ESC[<row>;<column>f` — то же самое (HVP)
			//		`ESC[<n>G` — курсор в столбец n (Cursor Horizontal Absolute, CHA)
			//		`ESC[<n>d` — курсор в строку n (Cursor Vertical Absolute, VPA)
			//Сохранение/восстановление
			//		`ESC[s` — сохранить позицию курсора (Save Cursor, DECSC)
			//		`ESC[u` — восстановить позицию курсора (Restore Cursor, DECRC)
			//		`ESC[6n` — запросить позицию курсора (DSR) → ответ: ESC[<row>;<column>R
			//Очистка экрана
			//		`ESC[0J` — очистить от курсора до конца экрана (ED)
			//		`ESC[1J` — очистить от начала до курсора
			//		`ESC[2J` — очистить весь экран
			//		`ESC[3J` — очистить весь экран и буфер прокрутки (xterm)
			//Очистка строки
			//		`ESC[0K` — очистить от курсора до конца строки (EL)
			//		`ESC[1K` — очистить от начала строки до курсора
			//		`ESC[2K` — очистить всю строку
			//256-цветная палитра
			//		`ESC[38;5;<n>m` — цвет текста (0-255)
			//		`ESC[48;5;<n>m` — цвет фона (0-255)
			//TrueColor (24-битный RGB)
			//		`ESC[38;2;<r>;<g>;<b>m` — RGB цвет текста
			//		`ESC[48;2;<r>;<g>;<b>m` — RGB цвет фона
			//6. Скроллинг (Scrolling)//@todo пока не требуется
			//		`ESC[<n>S` — прокрутить вверх на n строк (SU)
			//		`ESC[<n>T` — прокрутить вниз на n строк (SD)
			//		`ESC[<start>;<end>r` — установить область скроллинга (DECSTBM)
			//8. Вкладки (Tabs)//@todo пока не требуется
			//		`ESC[0g` — очистить табуляцию в текущей позиции (TBC)
			//		`ESC[3g` — очистить все табуляции (TBC)
			//		`ESC[I` — курсор к следующей табуляции (TAB)
			//9. Принтер//@todo пока не требуется
			//		`ESC[5i` — начать печать (включить принтер)
			//		`ESC[4i` — закончить печать (выключить принтер)
			//		`ESC[i` — распечатать экран
			//10. Состояние терминала (Device Status Reports)//@todo пока не требуется
			//		`ESC[5n` — запросить статус → ответ: ESC[0n (OK)
			//		`ESC[6n` — запросить позицию курсора → ответ: ESC[<r>;<c>R
			//		`ESC[c` — запросить идентификатор → ответ: ESC[?1;0c
			//		`ESC[0c` → ответ: ESC[?1;0c (VT100)
			//11. Прочие важные//@todo пока не требуется
			//		`ESC[c` — сброс терминала (RIS)
			if len(strRunes) <= i {
				return 1, nil
			}
			sl, offset := par.parseUintSlice(strRunes[i:])
			i += offset
			if len(sl) == 0 || len(strRunes) <= i {
				return 1, nil
			}
			if f, b := par.ansiMainList.GetFunction(strRunes[i]); b {
				if command, ok = f(sl); !ok {
					return 1, nil
				}

				return i + 1, command
			}
			if len(sl) == 1 {
				if f, b := par.ansiMainList.GetFunction('G'); b {
					if command, ok = f(sl); !ok {
						return 1, nil
					}

					return i, command
				}
			}

			return 1, nil
		}
	case ']':
		//Заголовок окна
		//		`ESC]0;<текст>\x07` — заголовок и иконки
		//		`ESC]1;<текст>\x07` — заголовок иконки
		//		`ESC]2;<текст>\x07` — заголовок окна
		//		`ESC]10;?` — запросить цвет текста
		//Палитра цветов
		//		`ESC]4;<n>;?` — запросить цвет палитры
		//		`ESC]4;<n>;rgb:<r>/<g>/<b>` — установить цвет палитры
		i++
		if len(strRunes) <= i {
			return 1, nil
		}
		sl, offset := par.parseUintSlice(strRunes[i:])
		i += offset
		if len(sl) == 0 || len(strRunes) <= i {
			return 1, nil
		}
		switch true {
		case strRunes[i] == '?':
			//Заголовок окна
			//		`ESC]10;?` — запросить цвет текста
			//Палитра цветов
			//		`ESC]4;<n>;?` — запросить цвет палитры
			if len(sl) < 2 || (sl[0] != 4 && sl[0] != 10) {
				return 1, nil
			}
			if f, b := par.ansiWindowAndPalettes.GetFunction(strRunes['?']); b {
				if command, ok = f(sl, ""); !ok {
					return 1, nil
				}
				return i + 1, command
			} else {
				return 1, nil
			}
		case len(strRunes) > i+3 && string(strRunes[i:i+3]) == "rgb":
			//Палитра цветов
			//		`ESC]4;<n>;rgb:<r>/<g>/<b>` — установить цвет палитры
			if len(sl) < 2 || sl[0] != 4 {
				return 1, nil
			}
			i += 3
			sl2, offset2 := par.parseUintSlice(strRunes[i:])
			i += offset2
			if f, b := par.ansiWindowAndPalettes.GetFunction(strRunes['s']); b {
				if command, ok = f(append(sl, sl2...), ""); !ok {
					return 1, nil
				}
				return i + 1, command
			} else {
				return 1, nil
			}
		default:
			//Заголовок окна
			//		`ESC]0;<текст>\x07` — заголовок и иконки
			//		`ESC]1;<текст>\x07` — заголовок иконки
			//		`ESC]2;<текст>\x07` — заголовок окна
			if len(sl) < 2 {
				return 1, nil
			}
			str, offset2 := par.parseString(strRunes[i:])
			if i+offset2 >= len(strRunes) {
				return 1, nil
			}
			i += offset2

			if f, b := par.ansiWindowAndPalettes.GetFunction(strRunes['t']); b {
				if command, ok = f(sl, str); !ok {
					return 1, nil
				}
				return i + 1, command
			} else {
				return 1, nil
			}

		}

	//case '('://@todo пока не требуется
	//7. Шрифты и клавиатура
	//		`ESC(0` — набор символов DEC
	//		`ESC(B` — набор символов ASCII

	//Сохранение/восстановление
	//		`ESC 7` — сохранить позицию (альтернатива)
	case '7':
		if f, b := par.ansiMainList.GetFunction('s'); b {
			if command, ok = f([]uint64{0}); !ok {
				return 1, nil
			}
			return i + 1, command
		} else {
			return 1, nil
		}
	//		ESC 8 — восстановить позицию (альтернатива)
	case '8':
		if f, b := par.ansiMainList.GetFunction('u'); b {
			if command, ok = f([]uint64{0}); !ok {
				return 1, nil
			}
			return i + 1, command
		} else {
			return 1, nil
		}
	default:
		//Сохранение/восстановление
		//		`ESC D` — индекс (IND) — курсор вниз со скроллингом//@todo пока не требуется
		//		`ESC M` — обратный индекс (RI) — курсор вверх со скроллингом//@todo пока не требуется
		//		`ESC#8` — заполнить экран символом 'E' (DECALN)//@todo пока не требуется
		//6. Скроллинг (Scrolling)
		//		`ESC D` — индекс (IND) — курсор вниз со скроллингом//@todo пока не требуется
		//		`ESC M` — обратный индекс (RI) — курсор вверх со скроллингом//@todo пока не требуется
		//8. Вкладки (Tabs)
		//		`ESC H` — установить табуляцию в текущей позиции (HTS)//@todo пока не требуется
		//11. Прочие важные
		//		`ESC#8` — заполнить экран символом 'E' (DECALN)//@todo пока не требуется
		//		`ESC>` — режим цифровой клавиатуры (DECKPNM)//@todo пока не требуется
		//		`ESC=` — режим цифровой клавиатуры (DECKPAM)//@todo пока не требуется
		//		`ESC Z` — запросить идентификатор → ответ: ESC[?1;0c (DECID)//@todo пока не требуется
	}

	return i, nil
}

func (par *parser) Parse(text string) []Command {
	var (
		textRunes   = []rune(text)
		command     Command
		commands    = make([]Command, 0)
		endString   int
		startString int
		offset      int
	)
	for i := 0; i < len(textRunes); {
		switch textRunes[i] {
		case '\n':
			command = par.getNewLineCommand()
		case '\r':
			command = par.getBeginningOfLineCommand()
		case '\033':
			if i+1 < len(textRunes) {
				if offset, command = par.enterInParseANSIEscapeSequences(textRunes[i:]); command != nil {
					i += offset
				} else {
					command = par.getTextCommand(textRunes, startString, endString)
					startString = i + 1
				}
			} else {
				command = par.getTextCommand(textRunes, startString, endString)
				startString = i + 1
			}
		default:
			if textRunes[i] < 0 {
				command = par.getTextCommand(textRunes, startString, endString)
				startString = i + 1
			}
		}
		if command != nil {
			if startString < endString {
				commands = append(commands, par.getTextCommand(textRunes, startString, endString), command)
			} else {
				commands = append(commands, command)
			}
			startString = i
			command = nil
		} else {
			i++
		}
		endString = i
	}
	if startString < len(textRunes) {
		commands = append(commands, par.getTextCommand(textRunes, startString, endString), par.getNewLineCommand())
	} else {
		commands = append(commands, par.getNewLineCommand())
	}

	return commands
}
