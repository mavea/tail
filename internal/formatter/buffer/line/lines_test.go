package line

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lines", func() {
	toStrings := func(items []Line) []string {
		result := make([]string, 0, len(items))
		for _, item := range items {
			Expect(item).NotTo(BeNil())
			result = append(result, item.String("", 0))
		}

		return result
	}

	setLine := func(lines Lines, id uint64, value string) {
		item := lines.Get(id)
		Expect(item).NotTo(BeNil())
		item.Set("", value, 0)
	}

	expectState := func(linesV Lines, expected []string) {
		Expect(toStrings(linesV.GetFullLines())).To(Equal(expected))

		x := linesV.LenHistory() - uint64(len(expected))
		for _, value := range expected {
			Expect(linesV.Get(x).String("", 0)).To(Equal(value))
			x++
		}
		Expect(linesV.Get(x)).To(BeNil())

		Expect(toStrings(linesV.GetLastLines(-1))).To(BeEmpty())
		Expect(toStrings(linesV.GetLastLines(int(linesV.LenHistory() + 5)))).To(Equal(expected))

		if len(expected) >= 2 {
			Expect(toStrings(linesV.GetLastLines(2))).To(Equal(expected[len(expected)-2:]))
		}

		last3Expected := expected
		if len(expected) > 3 {
			last3Expected = expected[len(expected)-3:]
		}
		Expect(toStrings(linesV.GetLastLines(3))).To(Equal(last3Expected))

		last4Expected := expected
		if len(expected) > 4 {
			last4Expected = expected[len(expected)-4:]
		}
		Expect(toStrings(linesV.GetLastLines(4))).To(Equal(last4Expected))
	}

	DescribeTable("инициализируется одной пустой строкой и держит границы через интерфейс",
		func(size uint64) {
			expectState(MakeLines(size), []string{""})
		},
		Entry("внесли 1 строку", uint64(1)),
		Entry("внесли 2 строки", uint64(2)),
		Entry("внесли 3 строки", uint64(3)),
		Entry("внесли 4 строки", uint64(4)),
	)

	It("корректно проходит 2-й и 3-й круг при размере 1", func() {
		lines := MakeLines(1)
		setLine(lines, 0, "L0")

		for i := 1; i <= 6; i++ {
			lines.Add(1)
			setLine(lines, lines.LenHistory()-1, fmt.Sprintf("L%d", i))
			expectState(lines, []string{fmt.Sprintf("L%d", i)})
		}
	})

	It("сохраняет последние элементы в правильном порядке на 2-м и 3-м круге", func() {
		lines := MakeLines(3)
		expected := []string{"L0"}
		setLine(lines, 0, expected[0])
		expectState(lines, expected)

		for i := 1; i <= 8; i++ {
			lines.Add(1)
			setLine(lines, lines.LenHistory()-1, fmt.Sprintf("L%d", i))

			expected = append(expected, fmt.Sprintf("L%d", i))
			if len(expected) > 3 {
				expected = expected[len(expected)-3:]
			}
			expectState(lines, expected)
		}
	})

	It("Add(count >= cap) полностью заменяет окно новыми строками", func() {
		lines := MakeLines(3)
		setLine(lines, 0, "A0")
		lines.Add(2)
		setLine(lines, 1, "A1")
		setLine(lines, 2, "A2")
		expectState(lines, []string{"A0", "A1", "A2"})

		lines.CleanPrefix(1)
		expectState(lines, []string{"A1", "A2"})

		lines.Add(5)
		expectState(lines, []string{""})
	})

	It("корректно чистит строку, постфикс и префикс после wrap-around", func() {
		lines := MakeLines(3)
		setLine(lines, 0, "L0")

		for i := 1; i <= 5; i++ {
			lines.Add(1)
			setLine(lines, lines.LenHistory()-1, fmt.Sprintf("L%d", i))
		}
		expectState(lines, []string{"L3", "L4", "L5"})

		lines.CleanString(1)
		expectState(lines, []string{"L3", "", "L5"})

		lines.CleanPostfix(1)
		expectState(lines, []string{"L3", ""})

		lines.CleanPrefix(1)
		expectState(lines, []string{""})
	})

	It("сбрасывается в одну пустую строку при полной очистке", func() {
		lines := MakeLines(3)
		setLine(lines, 0, "L0")
		lines.Add(2)
		setLine(lines, 1, "L1")
		setLine(lines, 2, "L2")
		lines.CleanPostfix(0)
		expectState(lines, []string{"L0"})

		setLine(lines, 0, "after-reset")
		expectState(lines, []string{"after-reset"})
		lines.CleanPrefix(9)
		expectState(lines, []string{""})
	})

	DescribeTable("буфер на 5 строк корректно теряет старые данные при большом потоке",
		func(written int) {
			linesV := MakeLines(5)
			for i := 0; i < written; i++ {
				setLine(linesV, linesV.LenHistory()-1, fmt.Sprintf("L%d", i))
				if i < written-1 {
					linesV.Add(1)
				}
			}

			window := 5
			if written < window {
				window = written
			}

			start := written - window
			expected := make([]string, 0, window)
			for i := start; i < written; i++ {
				expected = append(expected, fmt.Sprintf("L%d", i))
			}

			expectState(linesV, expected)
			Expect(toStrings(linesV.GetFullLines())).To(Equal(expected))

			last4Expected := expected
			if len(expected) > 4 {
				last4Expected = expected[len(expected)-4:]
			}
			Expect(toStrings(linesV.GetLastLines(4))).To(Equal(last4Expected))

			x := linesV.LenHistory() - uint64(len(expected))
			for _, value := range expected {
				Expect(linesV.Get(x).String("", 0)).To(Equal(value))
				x++
			}
			Expect(linesV.Get(x)).To(BeNil())
		},
		Entry("внесли 4 строки", 4),
		Entry("внесли 5 строк", 5),
		Entry("внесли 6 строк", 6),
		Entry("внесли 7 строк", 7),
		Entry("внесли 9 строк", 9),
		Entry("внесли 10 строк", 10),
		Entry("внесли 11 строк", 11),
		Entry("внесли 12 строк", 12),
	)
})
