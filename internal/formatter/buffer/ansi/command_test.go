package ansi

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command", func() {
	Context("CommandUint", func() {
		DescribeTable("Sunny",
			func(command *commandUint, wantType TypeCommand, wantVals []uint64) {
				Expect(command.GetType()).To(Equal(wantType))
				Expect(command.GetValsUint()).To(Equal(wantVals))
				Expect(command.GetValsString()).To(Equal(""))
			},
			Entry(
				"empty values",
				&commandUint{key: 'A', vals: []uint64{}},
				TypeCommand('A'),
				[]uint64{},
			),
			Entry(
				"single value",
				&commandUint{key: 'B', vals: []uint64{42}},
				TypeCommand('B'),
				[]uint64{42},
			),
			Entry(
				"multiple values",
				&commandUint{key: 'C', vals: []uint64{1, 2, 3}},
				TypeCommand('C'),
				[]uint64{1, 2, 3},
			),
		)
	})
	Context("CommandUintString", func() {
		DescribeTable("Sunny",
			func(command *commandUintString, wantType TypeCommand, wantVals []uint64, wantStr string) {
				Expect(command.GetType()).To(Equal(wantType))
				Expect(command.GetValsUint()).To(Equal(wantVals))
				Expect(command.GetValsString()).To(Equal(wantStr))
			},
			Entry(
				"empty values and string",
				&commandUintString{key: 'D', vals: []uint64{}, valsStr: ""},
				TypeCommand('D'),
				[]uint64{},
				"",
			),
			Entry(
				"single value and string",
				&commandUintString{key: 'E', vals: []uint64{99}, valsStr: "value"},
				TypeCommand('E'),
				[]uint64{99},
				"value",
			),
			Entry(
				"multiple values and string",
				&commandUintString{key: 'F', vals: []uint64{4, 5, 6}, valsStr: "multiple"},
				TypeCommand('F'),
				[]uint64{4, 5, 6},
				"multiple",
			),
		)
	})
})
