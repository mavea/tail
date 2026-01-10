package buffer

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TestCfg struct {
	ico  string
	name string
}

func (t *TestCfg) GetLengthLines() int {
	return 50
}

func (t *TestCfg) GetSizeBuffer() uint64 {
	return 100
}

func (t *TestCfg) GetProcessIcon() string {
	return t.ico
}

func (t *TestCfg) SetProcessIcon(ico string) {
	t.ico = ico
}

func (t *TestCfg) GetCountLines() uint64 {
	return 10
}
func (t *TestCfg) GetProcessName() string {
	return t.name
}
func (t *TestCfg) SetProcessName(str string) {
	t.name = str
}
func (t *TestCfg) GetTemplate() string {
	return "none"
}
func (t *TestCfg) GetIcon() string {
	return "none"
}

var _ = Describe("Modules", func() {
	var (
		b Buffer
	)
	BeforeEach(func() {
		b = New(&TestCfg{name: "test"})
		b.SetDefaultStyle("\033[0m", "\u001B[0m", "|")
	})
	Context("Buffer And Parse", func() {
		It("Sunny", func() {
			b.Add("12345")
			b.Add("67890")
			b.Add("abcde")

			Expect(b.GetLast(4)).To(Equal([]string{
				"\x1b[0m12345|",
				"\x1b[0m67890|",
				"\x1b[0mabcde|",
				"\x1b[0m|",
			}))
			b.Add("*\033[2Ag")
			Expect(b.GetLast(4)).To(Equal([]string{
				"\x1b[0m12345|",
				"\x1b[0m6g890|",
				"\x1b[0mabcde|",
				"\x1b[0m*|",
			}))
			//Позиционирование
			//		ESC[<n>A — курсор вверх на n строк (Cursor Up, CUU)
			b.Add("#-\033[1Ah")
			Expect(b.GetLast(4)).To(Equal([]string{
				"\x1b[0m12345|",
				"\x1b[0m6gh90|",
				"\x1b[0m#-cde|",
				"\x1b[0m*|",
			}))
			b.Add("+-\033[1Ahijklm")
			Expect(b.GetLast(4)).To(Equal([]string{
				"\x1b[0m12345|",
				"\x1b[0m6ghijklm|",
				"\x1b[0m+-cde|",
				"\x1b[0m*|",
			}))
			//		ESC[<n>B — курсор вниз на n строк (Cursor Down, CUD)
			b.Add("#-\033[1Bnl")
			Expect(b.GetLast(5)).To(Equal([]string{
				"\x1b[0m12345|",
				"\x1b[0m6ghijklm|",
				"\x1b[0m#-cde|",
				"\x1b[0m* nl|",
				"\x1b[0m|",
			}))
			b.Add("\u001B[1A\033[?1048h\u001B[?1049h")
			Expect(b.GetLast(5)).To(Equal([]string{
				"\x1b[0m12345|",
				"\x1b[0m6ghijklm|",
				"\x1b[0m#-cde|",
				"\x1b[0m* nl|",
				"\x1b[0m|",
			}))
			b.Add("\u001B[4A\033[?1050h")
			Expect(b.GetLast(5)).To(Equal([]string{
				"\x1b[0m[?1050h|",
				"\x1b[0m6ghijklm|",
				"\x1b[0m#-cde|",
				"\x1b[0m* nl|",
				"\x1b[0m|",
			}))
			b.Add("\x1b[48;2;1;10;100mrgb\x1b[47mbg\x1b[47;36mbgcol\x1b[36;47;2m*\u001B[0mClean")
			Expect(b.GetLast(5)).To(Equal([]string{
				"\x1b[0m[?1050h|",
				"\x1b[0m\x1b[48;2;1;10;100mrgb\x1b[47mbg\x1b[47;36mbgcol\x1b[47;36;2m*\x1b[0mClean|",
				"\x1b[0m#-cde|",
				"\x1b[0m* nl|",
				"\x1b[0m|",
			}))
			//		ESC[<n>C — курсор вправо на n символов (Cursor Forward, CUF)
			//		ESC[<n>D — курсор влево на n символов (Cursor Backward, CUB)
			b.Add("\033[6CZ\033[2DY")
			Expect(b.GetLast(5)).To(Equal([]string{
				"\x1b[0m[?1050h|",
				"\x1b[0m\x1b[48;2;1;10;100mrgb\x1b[47mbg\x1b[47;36mbgcol\x1b[47;36;2m*\x1b[0mClean|",
				"\x1b[0m#-cdeYZ|",
				"\x1b[0m* nl|",
				"\x1b[0m|",
			}))
			//		ESC[<n>E — курсор вниз на n строк, в начало строки (CNL)
			//		ESC[<n>F — курсор вверх на n строк, в начало строки (CPL)
			b.Add("---\033[2F0123456\033[3E4567890")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?1050h|",
				"\x1b[0m0123456\x1b[47;36mcol\x1b[47;36;2m*\x1b[0mClean|",
				"\x1b[0m#-cdeYZ|",
				"\x1b[0m---l|",
				"\x1b[0m4567890|",
				"\x1b[0m|",
			}))
			//Абсолютное позиционирование
			//		ESC[<row>;<column>H — переместить курсор (Cursor Position, CUP)
			//		ESC[<row>;<column>f — то же самое (HVP)
			b.Add("+\033[0;5H|poIuyt\033[3;2fDaetPlus")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?105|poIuyt|",
				"\x1b[0m0123456\x1b[47;36mcol\x1b[47;36;2m*\x1b[0mClean|",
				"\x1b[0m#-cdeYZ|",
				"\x1b[0m--DaetPlus|",
				"\x1b[0m4567890|",
				"\x1b[0m+|",
			}))

			//		ESC[<n>G — курсор в столбец n (Cursor Horizontal Absolute, CHA)
			//		ESC[<n> — курсор в столбец n (Cursor Horizontal Absolute, CHA)
			//		ESC[<n>d — курсор в строку n (Cursor Vertical Absolute, VPA)
			b.Add("_\033[0GA\033[2#\033[3dC")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?105|poIuyt|",
				"\x1b[0m0123456\x1b[47;36mcol\x1b[47;36;2m*\x1b[0mClean|",
				"\x1b[0m#-cdeYZ|",
				"\x1b[0m--DCetPlus|",
				"\x1b[0mA5#7890|",
				"\x1b[0m+|",
			}))

			b.Add("L")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?105|poIuyt|",
				"\x1b[0m0123456\x1b[47;36mcol\x1b[47;36;2m*\x1b[0mClean|",
				"\x1b[0m#-cdeYZ|",
				"\x1b[0m--DCetPlus|",
				"\x1b[0mL5#7890|",
				"\x1b[0m+|",
			}))
			//10. Состояние терминала (Device Status Reports)
			//		ESC[5n — запросить статус → ответ: ESC[0n (OK)
			//		ESC[6n — запросить позицию курсора → ответ: ESC[<r>;<c>R
			//Сохранение/восстановление
			//		ESC[6n — запросить позицию курсора (DSR) → ответ: ESC[<row>;<column>R

			//Сохранение/восстановление
			//		ESC[s — сохранить позицию курсора (Save Cursor, DECSC)
			//		ESC[u — восстановить позицию курсора (Restore Cursor, DECRC)
			b.Add("\033[1F#\033[0;2H$\033[s@\033[2;2H}\033[u*")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?$*5|poIuyt|",
				"\x1b[0m0123456\x1b[47;36mcol\x1b[47;36;2m*\x1b[0mClean|",
				"\x1b[0m#-}deYZ|",
				"\x1b[0m--DCetPlus|",
				"\x1b[0m#5#7890|",
				"\x1b[0m+|",
			}))

			//Очистка строки
			//		ESC[0K — очистить от курсора до конца строки (EL)
			b.Add("\033[1;10H{\033[0K}")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?$*5|poIuyt|",
				"\x1b[0m0123456\x1b[47;36mcol\x1b[0m{}|",
				"\x1b[0m#-}deYZ|",
				"\x1b[0m--DCetPlus|",
				"\x1b[0m#5#7890|",
				"\x1b[0m+|",
			}))

			// позиционирование
			b.Add("\033[1;10H12!")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?$*5|poIuyt|",
				"\x1b[0m0123456\x1b[47;36mcol\x1b[0m12!|",
				"\x1b[0m#-}deYZ|",
				"\x1b[0m--DCetPlus|",
				"\x1b[0m#5#7890|",
				"\x1b[0m+|",
			}))
			b.Add("\033[1;10H!")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?$*5|poIuyt|",
				"\x1b[0m0123456\x1b[47;36mcol\x1b[0m!2!|",
				"\x1b[0m#-}deYZ|",
				"\x1b[0m--DCetPlus|",
				"\x1b[0m#5#7890|",
				"\x1b[0m+|",
			}))
			//		ESC[1K — очистить от начала строки до курсора
			b.Add("\033[1;10H\033[1K&")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?$*5|poIuyt|",
				"\x1b[0m&2!|",
				"\x1b[0m#-}deYZ|",
				"\x1b[0m--DCetPlus|",
				"\x1b[0m#5#7890|",
				"\x1b[0m+|",
			}))
			//		ESC[2K — очистить всю строку
			b.Add("\033[1;10H\033[2KNewText")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?$*5|poIuyt|",
				"\x1b[0mNewText|",
				"\x1b[0m#-}deYZ|",
				"\x1b[0m--DCetPlus|",
				"\x1b[0m#5#7890|",
				"\x1b[0m+|",
			}))

			//Очистка экрана
			//		ESC[0J — очистить от курсора до конца экрана (ED)
			b.Add("\033[1;5H\033[0J->")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?$*5|poIuyt|",
				"\x1b[0mNewTe->|",
				"\x1b[0m|",
			}))
			b.Add("\033[1;5H\033[0J->\033[0;0H")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m[?$*5|poIuyt|",
				"\x1b[0mNewTe->|",
			}))
			//		ESC[1J — очистить от начала до курсора
			b.Add("\033[1;3H\033[1J<-")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m<-->|",
				"\x1b[0m|",
			}))
			//		ESC[2J — очистить весь экран
			b.Add("123456\033[0;3H\033[2J{}")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m{}|",
				"\x1b[0m|",
			}))
			//		ESC[3J — очистить весь экран и буфер прокрутки (xterm)
			b.Add("123456\033[0;1H\033[3J><")
			Expect(b.GetLast(10)).To(Equal([]string{
				"\x1b[0m><|",
				"\x1b[0m|",
			}))

		})
	})

})
