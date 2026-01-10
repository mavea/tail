package parser

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"tail/internal/buffer/ansi"
)

type testsCom struct {
	key     ansi.TypeCommand
	vals    []uint64
	valsStr string
}

func (c *testsCom) GetType() ansi.TypeCommand {
	return c.key
}
func (c *testsCom) GetValsUint() []uint64 {
	return c.vals
}
func (c *testsCom) GetValsString() string {
	return c.valsStr
}

type MockWindowAndPalettes struct {
	mock.Mock
}

func (m *MockWindowAndPalettes) GetFunction(char rune) (func([]uint64, string) (ansi.Command, bool), bool) {
	args := m.Called(char)
	if f, ok := args.Get(0).(func([]uint64, string) (ansi.Command, bool)); ok {
		return f, args.Bool(1)
	}

	return nil, args.Bool(1)
}

type MockMainList struct {
	mock.Mock
}

func (m *MockMainList) GetFunction(char rune) (func([]uint64) (ansi.Command, bool), bool) {
	args := m.Called(char)
	if f, ok := args.Get(0).(func([]uint64) (ansi.Command, bool)); ok {
		return f, args.Bool(1)
	}

	return nil, args.Bool(1)
}

var _ = Describe("Parser", func() {
	var (
		mockWindowAndPalettesRepo *MockWindowAndPalettes
		mockMainListRepo          *MockMainList
		par                       *parser
	)
	BeforeEach(func() {
		mockWindowAndPalettesRepo = new(MockWindowAndPalettes)
		mockMainListRepo = new(MockMainList)
		par = &parser{
			ansiWindowAndPalettes: mockWindowAndPalettesRepo,
			ansiMainList:          mockMainListRepo,
		}
	})
	Context("parseUintSlice", func() {
		DescribeTable("Sunny",
			func(str string, rez []uint64, offset int) {
				re, i := par.parseUintSlice([]rune(str))
				Expect(re).To(Equal(rez))
				Expect(i).To(Equal(offset))
			},
			Entry("D 123", "D 123", []uint64{0}, 0),
			Entry("1D 123", "1D 123", []uint64{1}, 1),
			Entry("10D 123", "10D 123", []uint64{10}, 2),
			Entry("19D 123", "19D 123", []uint64{19}, 2),
			Entry(";D 123", ";D 123", []uint64{0, 0}, 1),
			Entry("1;D 123", "1;D 123", []uint64{1, 0}, 2),
			Entry(";1D 123", ";1D 123", []uint64{0, 1}, 2),
			Entry("1;1D 123", "1;1D 123", []uint64{1, 1}, 3),
			Entry(";;;1D 123", ";;;1D 123", []uint64{0, 0, 0, 1}, 4),
			Entry(";;2;1D 123", ";;2;1D 123", []uint64{0, 0, 2, 1}, 5),
			Entry("D", "D", []uint64{0}, 0),
			Entry("zero string", "", []uint64{}, 0),
			Entry("nil", nil, []uint64{}, 0),
		)
	})
	Context("parseString", func() {
		DescribeTable("Sunny",
			func(r []rune, str string, offset int) {
				st, i := par.parseString(r)
				Expect(st).To(Equal(str))
				Expect(i).To(Equal(offset))
			},
			Entry("nil", nil, "", 0),
			Entry("empty", []rune{}, "", 0),
			Entry("escape", []rune{'\033'}, "", 0),
			Entry("h escape", []rune{'h', '\033'}, "", 0),
			Entry("h", []rune{'h'}, "", 0),
			Entry("hh", []rune{'h', 'h'}, "", 0),
			Entry("string hhh", []rune{'h', 'h', '\x07'}, "hh", 3),
		)
	})
	Context("getTextCommand", func() {
		DescribeTable("Sunny",
			func(str []rune, start, end int, orin *testsCom) {
				mockWindowAndPalettesRepo.On("GetFunction", rune(5)).Return(
					func(garb []uint64, s string) (ansi.Command, bool) {
						return &testsCom{
							key:     ansi.KeyText,
							vals:    nil,
							valsStr: s,
						}, true
					}, true).Maybe()

				com := par.getTextCommand(str, start, end)
				if orin == nil {
					Expect(com).To(BeNil())
				} else {
					Expect(com).To(Not(BeNil()))
					Expect(com.GetType()).To(Equal(orin.GetType()))
					Expect(com.GetValsUint()).To(Equal(orin.GetValsUint()))
					Expect(com.GetValsString()).To(Equal(orin.GetValsString()))
				}
			},
			Entry("nil", nil, 0, 5, nil),
			Entry("empty string", []rune{}, 0, 0, nil),
			Entry("empty string 2", []rune(""), 0, 0, nil),
			Entry("empty string and error index", []rune(""), 0, 5, nil),
			Entry("short string and error index", []rune("b"), 0, 5, nil),
			Entry("string and error index", []rune("str"), 0, 5, nil),
			Entry("string and error index", []rune("str"), 0, 3, &testsCom{key: ansi.KeyText, vals: nil, valsStr: "str"}),
			Entry("string", []rune("full string in text"), 0, 6, &testsCom{key: ansi.KeyText, vals: nil, valsStr: "full s"}),
			Entry("full string", []rune("full string in text"), 0, 19, &testsCom{key: ansi.KeyText, vals: nil, valsStr: "full string in text"}),
		)
	})
	Context("getBeginningOfLineCommand", func() {
		It("Sunny", func() {
			mockMainListRepo.On("GetFunction", rune('G')).Return(
				func(garb []uint64) (ansi.Command, bool) {
					return &testsCom{
						key:  ansi.KeyCursorSetX,
						vals: nil,
					}, true
				}, true).Maybe()

			com := par.getBeginningOfLineCommand()
			Expect(com).To(Not(BeNil()))
			Expect(com.GetType()).To(Equal(ansi.KeyCursorSetX))
			Expect(com.GetValsUint()).To(BeNil())
		})
	})
	Context("getNewLineCommand", func() {
		It("Sunny", func() {
			mockMainListRepo.On("GetFunction", rune('E')).Return(
				func(garb []uint64) (ansi.Command, bool) {
					return &testsCom{
						key:  ansi.KeyCursorDownLeft,
						vals: nil,
					}, true
				}, true).Maybe()

			com := par.getNewLineCommand()
			Expect(com).To(Not(BeNil()))
			Expect(com.GetType()).To(Equal(ansi.KeyCursorDownLeft))
			Expect(com.GetValsUint()).To(BeNil())
		})
	})
	Context("enterInParseANSIEscapeSequences", func() {
		DescribeTable("Sunny",
			func(str []rune, offset int, orin *testsCom) {
				mockMainListRepo.On("GetFunction", rune('E')).Return(
					func(garb []uint64) (ansi.Command, bool) {
						return &testsCom{
							key:  ansi.KeyCursorDownLeft,
							vals: garb,
						}, true
					}, true).Maybe()
				mockWindowAndPalettesRepo.On("GetFunction", rune(5)).Return(
					func(garb []uint64, s string) (ansi.Command, bool) {
						return &testsCom{
							key:     ansi.KeyText,
							vals:    garb,
							valsStr: s,
						}, true
					}, true).Maybe()

				i, com := par.enterInParseANSIEscapeSequences(str)
				if orin == nil {
					Expect(com).To(BeNil())
				} else {
					Expect(com).To(Not(BeNil()))
					Expect(com.GetType()).To(Equal(orin.GetType()))
					Expect(com.GetValsUint()).To(Equal(orin.GetValsUint()))
				}
				Expect(i).To(Equal(offset))
			},
			Entry("\\033[3E", []rune("\033[3E"), 4, &testsCom{key: ansi.KeyCursorDownLeft, vals: []uint64{3}}),
			Entry("\\033[3E", []rune("\033[E"), 3, &testsCom{key: ansi.KeyCursorDownLeft, vals: []uint64{0}}),
			Entry("\\033[33", []rune("\033[33"), 1, nil),
		)
	})
	Context("Parse", func() {
		DescribeTable("Sunny",
			func(str string, orin []*testsCom) {
				mockMainListRepo.On("GetFunction", rune('E')).Return(
					func(garb []uint64) (ansi.Command, bool) {
						return &testsCom{
							key:  ansi.KeyCursorDownLeft,
							vals: garb,
						}, true
					}, true).Maybe()
				mockMainListRepo.On("GetFunction", rune('G')).Return(
					func(garb []uint64) (ansi.Command, bool) {
						return &testsCom{
							key:  ansi.KeyCursorSetX,
							vals: nil,
						}, true
					}, true).Maybe()
				mockWindowAndPalettesRepo.On("GetFunction", rune(5)).Return(
					func(garb []uint64, s string) (ansi.Command, bool) {
						return &testsCom{
							key:     ansi.KeyText,
							vals:    garb,
							valsStr: s,
						}, true
					}, true).Maybe()

				commands := par.Parse(str)
				if orin == nil {
					Expect(commands).To(BeNil())
				} else {
					for i, com := range commands {
						if i < len(orin) {
							if orin[i] == nil {
								Expect(com).To(BeNil())
							} else {
								Expect(com).To(Not(BeNil()))
								Expect(com.GetType()).To(Equal(orin[i].GetType()))
								Expect(com.GetValsUint()).To(Equal(orin[i].GetValsUint()))
								Expect(com.GetValsString()).To(Equal(orin[i].GetValsString()))
							}
						} else {
							Expect(com).To(BeNil())
						}
					}
				}
			},
			Entry("\\033[3Etext", "\033[3Etext", []*testsCom{
				{key: ansi.KeyCursorDownLeft, vals: []uint64{3}},
				{key: ansi.KeyText, valsStr: "text"},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{1}},
			}),
			Entry("\\033[3Etext\\033[80Etext2", "\u001B[3Etext\u001B[80Etext2", []*testsCom{
				{key: ansi.KeyCursorDownLeft, vals: []uint64{3}},
				{key: ansi.KeyText, valsStr: "text"},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{80}},
				{key: ansi.KeyText, valsStr: "text2"},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{1}},
			}),
			Entry("prefix\\033[3Etext\\033[80Etext2", "prefix\u001B[3Etext\u001B[80Etext2", []*testsCom{
				{key: ansi.KeyText, valsStr: "prefix"},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{3}},
				{key: ansi.KeyText, valsStr: "text"},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{80}},
				{key: ansi.KeyText, valsStr: "text2"},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{1}},
			}),
			Entry("prefix\\033[3E", "prefix\u001B[3E", []*testsCom{
				{key: ansi.KeyText, valsStr: "prefix"},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{3}},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{1}},
			}),
			Entry("empty", "", []*testsCom{
				{key: ansi.KeyCursorDownLeft, vals: []uint64{1}},
			}),
			Entry("\\u001B!!", "\u001B!!", []*testsCom{
				{key: ansi.KeyText, valsStr: "!!"},
				{key: ansi.KeyCursorDownLeft, vals: []uint64{1}},
			}),
		)
	})
})
