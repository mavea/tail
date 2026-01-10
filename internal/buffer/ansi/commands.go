package ansi

// DECPrivateModes
type DECPrivateModes []func(a []uint64) (Command, bool)

func (d DECPrivateModes) getKey(char rune) rune {
	switch char {
	case 'h':
		return 1
	case 'l':
		return 2
	case 'c':
		return 3
	}

	return 0
}
func (d DECPrivateModes) GetFunction(char rune) (func([]uint64) (Command, bool), bool) {
	if k := d.getKey(char); k != 0 {
		if d := d[k]; d != nil {
			return d, true
		}
	}

	return nil, false
}

func GetActionsDECPrivateModes() DECPrivateModes {
	return DECPrivateModes{
		nil,
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
		//Мышь
		//		`ESC[?1000h` — включить мышь (X10)
		//		`ESC[?1002h` — включить мышь (для drag)
		//		`ESC[?1003h` — включить мышь (все события)
		//		`ESC[?1006h` — SGR-формат координат мыши
		//Фокус
		//		`ESC[?1004h` — отслеживание фокуса
		func(val []uint64) (Command, bool) {
			if len(val) == 1 {
				switch val[0] {
				case 1, 3, 5, 6, 7, 8, 12, 25, 47, 1047, 1048, 1049: //мышь и фокус игнорим
					return &commandUint{key: KeyOtherDo,
						vals: val,
					}, true
				default:
					return nil, false
				}
			}
			return nil, false
		},
		//Сброс режимов (Reset Mode)
		//		`ESC[?<n>l` (буква L)
		//				(Те же номера, что и выше)
		func(val []uint64) (Command, bool) {
			if len(val) == 1 {
				switch val[0] {
				case 1, 3, 5, 6, 7, 8, 12, 25, 47, 1047, 1048, 1049: //мышь и фокус игнорим
					return &commandUint{key: KeyOtherRevert,
						vals: val,
					}, true
				default:
					return nil, false
				}
			}
			return nil, false
		},
		//7. Шрифты и клавиатура
		//		`ESC[?1c` — ответ: ESC[?1;0c (базовый VT100)
		func(val []uint64) (Command, bool) {
			if (len(val) == 1 || len(val) == 2) && val[0] == 1 {
				return &commandUint{key: KeyVT100,
					vals: []uint64{1, 0},
				}, true
			}
			return nil, false
		},
	}
}

// SecondaryDA
type SecondaryDA []func(a []uint64) (Command, bool)

func (s SecondaryDA) getKey(char rune) rune {
	switch char {
	case 'c':
		return 1
	case 'm':
		return 2
	}

	return 0
}
func (s SecondaryDA) GetFunction(char rune) (func([]uint64) (Command, bool), bool) {
	if k := s.getKey(char); k != 0 {
		if d := s[k]; d != nil {
			return d, true
		}
	}

	return nil, false
}
func GetActionsSecondaryDA() SecondaryDA {
	//7. Шрифты и клавиатура
	//		`ESC[>c` — запросить идентификатор терминала
	return SecondaryDA{
		nil,
		func(val []uint64) (Command, bool) {
			if (len(val) == 1 || len(val) == 2) && val[0] == 1 {
				return &commandUint{key: KeyTerminalID,
					vals: []uint64{1, 0},
				}, true
			}
			return nil, false
		},
		//7. Шрифты и клавиатура
		//		`ESC[>0;256m` — модификаторы клавиш (xterm)
		func(val []uint64) (Command, bool) {
			if (len(val) == 1 || len(val) == 2) && val[0] == 1 {
				return &commandUint{key: KeyKeysModificator,
					vals: []uint64{1, 0},
				}, true
			}
			return nil, false
		},
	}
}

// WindowAndPalettes
type WindowAndPalettes []func([]uint64, string) (Command, bool)

func (w WindowAndPalettes) getKey(char rune) rune {
	switch char {
	case '?':
		return 1
	case 's':
		return 2
	case 't':
		return 3
	case 5:
		return 4
	}

	return 0
}
func (w WindowAndPalettes) GetFunction(char rune) (func([]uint64, string) (Command, bool), bool) {
	if k := w.getKey(char); k != 0 {
		if d := w[k]; d != nil {
			return d, true
		}
	}

	return nil, false
}
func GetActionsWindowAndPalettes() WindowAndPalettes {
	//Заголовок окна
	//		`ESC]0;<текст>\x07` — заголовок и иконки
	//		`ESC]1;<текст>\x07` — заголовок иконки
	//		`ESC]2;<текст>\x07` — заголовок окна
	//		`ESC]10;?` — запросить цвет текста
	//Палитра цветов
	//		`ESC]4;<n>;?` — запросить цвет палитры
	//		`ESC]4;<n>;rgb:<r>/<g>/<b>` — установить цвет палитры
	return WindowAndPalettes{
		nil,
		func(val []uint64, str string) (Command, bool) {
			//Заголовок окна
			//		`ESC]10;?` — запросить цвет текста
			//Палитра цветов
			//		`ESC]4;<n>;?` — запросить цвет палитры
			if len(val) == 2 {
				return &commandUint{key: KeyGetPalettes,
					vals: val,
				}, true
			}
			return nil, false
		},
		func(val []uint64, str string) (Command, bool) {
			//		`ESC]4;<n>;rgb:<r>/<g>/<b>` — установить цвет палитры
			if len(val) == 2 {
				return &commandUint{key: KeySetPalettes,
					vals: val,
				}, true
			}
			return nil, false
		},
		func(val []uint64, str string) (Command, bool) {
			//Заголовок окна
			//		`ESC]0;<текст>\x07` — заголовок и иконки
			//		`ESC]1;<текст>\x07` — заголовок иконки
			//		`ESC]2;<текст>\x07` — заголовок окна
			if len(val) == 2 {
				return &commandUintString{key: KeyTitle,
					vals:    val,
					valsStr: str,
				}, true
			}
			return nil, false
		},
		//просто текст
		func(val []uint64, str string) (Command, bool) {
			return &commandUintString{key: KeyText,
				valsStr: str,
			}, true
		},
	}
}

// FromMainList
type MainList []func(a []uint64) (Command, bool)

func (m MainList) getKey(char rune) rune {
	switch char {
	case 'A':
		return 1
	case 'B':
		return 2
	case 'C':
		return 3
	case 'D':
		return 4
	case 'E':
		return 5
	case 'F':
		return 6
	case 'H', 'f':
		return 7
	case 'G':
		return 8
	case 'd':
		return 9
	case 'n':
		return 10
	case 's':
		return 11
	case 'u':
		return 12
	case 'J':
		return 13
	case 'K':
		return 14
	case 'm':
		return 15
	}

	return 0
}

func (m MainList) GetFunction(char rune) (func([]uint64) (Command, bool), bool) {
	if k := m.getKey(char); k != 0 {
		if d := m[k]; d != nil {
			return d, true
		}
	}

	return nil, false
}
func GetActionsMainList() MainList {
	return MainList{
		nil,
		//Позиционирование
		//		`ESC[<n>A` — курсор вверх на n строк (Cursor Up, CUU)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{1}
			}
			return &commandUint{key: KeyCursorUp,
				vals: val,
			}, true
		},
		//		`ESC[<n>B` — курсор вниз на n строк (Cursor Down, CUD)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{1}
			}
			return &commandUint{key: KeyCursorDown,
				vals: val,
			}, true
		},
		//		`ESC[<n>C` — курсор вправо на n символов (Cursor Forward, CUF)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{1}
			}
			return &commandUint{key: KeyCursorRight,
				vals: val,
			}, true
		},
		//		`ESC[<n>D` — курсор влево на n символов (Cursor Backward, CUB)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{1}
			}
			return &commandUint{key: KeyCursorLeft,
				vals: val,
			}, true
		},
		//		`ESC[<n>E` — курсор вниз на n строк, в начало строки (CNL)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{1}
			}
			return &commandUint{key: KeyCursorDownLeft,
				vals: val,
			}, true
		},
		//		`ESC[<n>F` — курсор вверх на n строк, в начало строки (CPL)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{1}
			}
			return &commandUint{key: KeyCursorUpLeft,
				vals: val,
			}, true
		},
		//Абсолютное позиционирование
		//		`ESC[<row>;<column>H` — переместить курсор (Cursor Position, CUP)
		//		`ESC[<row>;<column>f` — то же самое (HVP)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{0, 0}
			}
			if len(val) == 1 {
				val = []uint64{val[0], 0}
			}
			return &commandUint{key: KeySetCursorXY,
				vals: val,
			}, true
		},
		//		`ESC[<n>G` — курсор в столбец n (Cursor Horizontal Absolute, CHA)
		//		`ESC[<n>` — курсор в столбец n (Cursor Horizontal Absolute, CHA)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{0}
			}
			return &commandUint{key: KeyCursorSetX,
				vals: val,
			}, true
		},
		//		`ESC[<n>d` — курсор в строку n (Cursor Vertical Absolute, VPA)
		func(val []uint64) (Command, bool) {
			if len(val) == 0 {
				val = []uint64{0}
			}
			return &commandUint{key: KeyCursorSetY,
				vals: val,
			}, true
		}, //9

		//10. Состояние терминала (Device Status Reports)
		//		`ESC[5n` — запросить статус → ответ: ESC[0n (OK)
		//		`ESC[6n` — запросить позицию курсора → ответ: ESC[<r>;<c>R
		//Сохранение/восстановление
		//		`ESC[6n` — запросить позицию курсора (DSR) → ответ: ESC[<row>;<column>R
		func(val []uint64) (Command, bool) {
			if len(val) != 1 || (val[0] != 6 && val[0] != 5) {
				return nil, false
			}
			return &commandUint{key: KeyGetStatusOrCursor,
				vals: val,
			}, true
		}, //10

		//Сохранение/восстановление
		//		`ESC[s` — сохранить позицию курсора (Save Cursor, DECSC)
		func(val []uint64) (Command, bool) {
			if len(val) != 1 || val[0] != 0 {
				return nil, false
			}
			return &commandUint{key: KeyCursorSave,
				vals: val,
			}, true
		}, //11
		//		`ESC[u` — восстановить позицию курсора (Restore Cursor, DECRC)
		func(val []uint64) (Command, bool) {
			if len(val) != 1 || val[0] != 0 {
				return nil, false
			}
			return &commandUint{key: KeyCursorLoad,
				vals: val,
			}, true
		}, //12
		//Очистка экрана
		//		`ESC[0J` — очистить от курсора до конца экрана (ED)
		//		`ESC[1J` — очистить от начала до курсора
		//		`ESC[2J` — очистить весь экран
		//		`ESC[3J` — очистить весь экран и буфер прокрутки (xterm)
		func(val []uint64) (Command, bool) {
			if len(val) != 1 || val[0] > 3 {
				return nil, false
			}
			return &commandUint{key: KeyClearScreen,
				vals: val,
			}, true
		}, //13
		//Очистка строки
		//		`ESC[0K` — очистить от курсора до конца строки (EL)
		//		`ESC[1K` — очистить от начала строки до курсора
		//		`ESC[2K` — очистить всю строку
		func(val []uint64) (Command, bool) {
			if len(val) != 1 || val[0] > 2 {
				return nil, false
			}
			return &commandUint{key: KeyClearLine,
				vals: val,
			}, true
		}, //14
		//Цвета и стили текста
		// 			`ESC[0m` # сброс всех атрибутов
		// 			`ESC[1m` # жирный (bright)
		// 			`ESC[2m` # тусклый (dim)
		// 			`ESC[3m` # курсив
		// 			`ESC[4m` # подчеркивание
		// 			`ESC[5m` # мигание
		// 			`ESC[7m` # инверсные цвета
		// 			`ESC[8m` # скрытый текст
		// 			`ESC[9m` # зачеркнутый

		//Цвета текста (основные)
		// 			`ESC[30m` # черный
		// 			`ESC[31m` # красный
		// 			`ESC[32m` # зеленый
		// 			`ESC[33m` # желтый
		// 			`ESC[34m` # синий
		// 			`ESC[35m` # пурпурный
		// 			`ESC[36m` # голубой
		// 			`ESC[37m` # белый
		// 			`ESC[90m` # ярко-черный (серый)
		// 			`ESC[91m` # ярко-красный
		// 			`ESC[92m` # ярко-зеленый
		// 			`ESC[93m` # ярко-желтый
		// 			`ESC[94m` # ярко-синий
		// 			`ESC[95m` # ярко-пурпурный
		// 			`ESC[96m` # ярко-голубой
		// 			`ESC[97m` # ярко-белый

		//Цвета фона
		// 			`ESC[40m` # черный
		// 			`ESC[41m` # красный
		// 			`ESC[42m` # зеленый
		// 			`ESC[43m` # желтый
		// 			`ESC[44m` # синий
		// 			`ESC[45m` # пурпурный
		// 			`ESC[46m` # голубой
		// 			`ESC[47m` # белый
		// 			`ESC[100m` # ярко-черный (серый)
		// 			`ESC[101m` # ярко-красный
		// 			`ESC[102m` # ярко-зеленый
		// 			`ESC[103m` # ярко-желтый
		// 			`ESC[104m` # ярко-синий
		// 			`ESC[105m` # ярко-пурпурный
		// 			`ESC[106m` # ярко-голубой
		// 			`ESC[107m` # ярко-белый
		//256-цветная палитра
		//		`ESC[38;5;<n>m` — цвет текста (0-255)
		//		`ESC[48;5;<n>m` — цвет фона (0-255)
		//TrueColor (24-битный RGB)
		//		`ESC[38;2;<r>;<g>;<b>m` — RGB цвет текста
		//		`ESC[48;2;<r>;<g>;<b>m` — RGB цвет фона
		func(val []uint64) (Command, bool) {

			if len(val) < 1 {
				return nil, false
			}
			if val[0] == 38 || val[0] == 48 {
				switch true {
				case len(val) < 4 && len(val) >= 2 && val[1] == 5 && (len(val) == 2 || val[2] < 255):
				case len(val) < 6 && len(val) >= 2 && val[1] == 2 &&
					(len(val) < 3 || val[2] < 255) &&
					(len(val) < 4 || val[3] < 255) &&
					(len(val) < 5 || val[4] < 255):
				default:
					return nil, false
				}
			}
			return &commandUint{key: KeyColor,
				vals: val,
			}, true
		}, //15
	}
}

type OtherList []func(a []uint64) (Command, bool)

func (o OtherList) GetFunction(char rune) (func([]uint64) (Command, bool), bool) {
	return nil, false
}

func GetActionsOtherList() OtherList {
	return OtherList{
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
		//		`ESC[c` — запросить идентификатор → ответ: ESC[?1;0c
		//		`ESC[0c` → ответ: ESC[?1;0c (VT100)
		//11. Прочие важные//@todo пока не требуется
		//		`ESC[c` — сброс терминала (RIS)
	}
}
func GetLists() (
	DECPrivateModes,
	SecondaryDA,
	WindowAndPalettes,
	MainList,
	OtherList,
) {
	return GetActionsDECPrivateModes(),
		GetActionsSecondaryDA(),
		GetActionsWindowAndPalettes(),
		GetActionsMainList(),
		GetActionsOtherList()
}
