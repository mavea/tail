package line

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Line", func() {
	Context("String", func() {
		DescribeTable("Sunny",
			func(initial line, clean string, length int, expected string) {
				Expect(initial.String(clean, length)).To(Equal(expected))
			},
			Entry(
				"empty line",
				line{},
				" ",
				0,
				"",
			),
			Entry(
				"no styles, single content",
				line{"", "content"},
				" ",
				0,
				"content",
			),
			Entry(
				"no styles, multiple entries",
				line{"", "content", "", "content2"},
				" ",
				0,
				"contentcontent2",
			),
			Entry(
				"multiple entries, mixed styles",
				line{"styleA", "content1", "styleB", "content2"},
				" ",
				0,
				"styleAcontent1styleBcontent2",
			),
			Entry(
				"multiple entries, mixed styles and length",
				line{"styleA", "content1", "styleA", "content2"},
				" ",
				8,
				"styleAcontent1 .",
			),
			Entry(
				"multiple entries, mixed styles and length 2",
				line{"styleA", "content1", "styleB", "content2"},
				" ",
				10,
				"styleAcontent1styleBco .",
			),
			Entry(
				"same style reused",
				line{"styleA", "content1", "styleA", "content2"},
				" ",
				0,
				"styleAcontent1content2",
			),
			Entry(
				"same style reused and length",
				line{"styleA", "content1", "styleA", "content2"},
				" ",
				10,
				"styleAcontent1co .",
			),
		)
	})
	Context("skipToX", func() {
		DescribeTable("Sunny",
			func(initial line, x int, expectedY int, expectedOffset int) {
				y, offset := initial.skipToX(x)
				Expect(offset).To(Equal(expectedOffset))
				Expect(y).To(Equal(expectedY))
			},
			Entry(
				"empty line",
				line{},
				5,
				1,
				-5,
			),
			Entry(
				"skipping within bounds",
				line{"styleA", "12345", "styleB", "123"},
				5,
				1,
				5,
			),
			Entry(
				"skipping with offset",
				line{"styleA", "12345", "styleB", "123"},
				4,
				1,
				4,
			),
			Entry(
				"skipping with offset2",
				line{"styleA", "12345", "styleB", "123"},
				7,
				3,
				2,
			),
			Entry(
				"skipping out of bounds",
				line{"styleA", "12345", "styleB", "123"},
				50,
				5,
				-42,
			),
			Entry(
				"negative skipping",
				line{"styleA", "12345", "styleB", "123"},
				-1,
				1,
				0,
			),
			Entry(
				"zero skipping",
				line{"styleA", "12345", "styleB", "123"},
				0,
				1,
				0,
			),
		)
	})
	Context("skipToXFromY", func() {
		DescribeTable("Sunny",
			func(initial line, x, y int, expectedY int, expectedOffset int) {
				y, offset := initial.skipToXFromY(x, y)
				Expect(offset).To(Equal(expectedOffset))
				Expect(y).To(Equal(expectedY))
			},
			Entry(
				"empty line",
				line{},
				5,
				1,
				1,
				-5,
			),
			Entry(
				"skipping within bounds",
				line{"styleA", "12345", "styleB", "123"},
				5,
				1,
				1,
				5,
			),
			Entry(
				"skipping with offset",
				line{"styleA", "12345", "styleB", "123"},
				4,
				1,
				1,
				4,
			),
			Entry(
				"skipping with offset2",
				line{"styleA", "12345", "styleB", "123"},
				7,
				1,
				3,
				2,
			),
			Entry(
				"skipping out of bounds",
				line{"styleA", "12345", "styleB", "123"},
				50,
				1,
				5,
				-42,
			),
			Entry(
				"negative skipping",
				line{"styleA", "12345", "styleB", "123"},
				-1,
				1,
				1,
				0,
			),
			Entry(
				"negative skipping",
				line{"styleA", "12345", "styleB", "123"},
				-1,
				2,
				2,
				0,
			),
			Entry(
				"zero skipping",
				line{"styleA", "12345", "styleB", "123"},
				0,
				2,
				2,
				0,
			),
			Entry(
				"zero skipping",
				line{"styleA", "12345", "styleB", "123"},
				3,
				2,
				2,
				3,
			),
			Entry(
				"zero skipping",
				line{"styleA", "12345", "styleB", "123", "styleB", "12"},
				4,
				3,
				5,
				1,
			),
			Entry(
				"zero skipping",
				line{"styleA", "12345", "styleB", "123", "styleB", "12"},
				9,
				3,
				7,
				-4,
			),
		)
	})
	Context("Set", func() {
		DescribeTable("Sunny",
			func(initial line, style, add string, x int, expected line, expectedX int) {
				offset := initial.Set(style, add, x)
				Expect(initial).To(Equal(expected))
				Expect(offset).To(Equal(expectedX))
			},
			Entry(
				"simple set with add",
				line{},
				"styleA",
				"new_text",
				0,
				line{"styleA", "new_text"},
				8,
			),
			Entry(
				"append to existing line",
				line{"styleA", "123"},
				"styleA",
				"456",
				3,
				line{"styleA", "123456"},
				6,
			),
			Entry(
				"append to line",
				line{"styleA", "123"},
				"styleA",
				"456",
				2,
				line{"styleA", "12456"},
				5,
			),
			Entry(
				"append to existing line different styles",
				line{"styleA", "123"},
				"styleB",
				"456",
				3,
				line{"styleA", "123", "styleB", "456"},
				6,
			),
			Entry(
				"different styles",
				line{"styleA", "123"},
				"styleB",
				"456",
				2,
				line{"styleA", "12", "styleB", "456"},
				5,
			),
			Entry(
				"different styles and space",
				line{"styleA", "123  "},
				"styleB",
				"456",
				5,
				line{"styleA", "123  ", "styleB", "456"},
				8,
			),
			Entry(
				"different styles and space",
				line{"styleA", "123"},
				"styleB",
				"456",
				5,
				line{"styleA", "123", "styleB", "  456"},
				8,
			),
		)
	})
	Context("CleanPrefix", func() {
		DescribeTable("Sunny",
			func(initial line, x int, expected line, expectedX int) {
				offset := initial.CleanPrefix(x)
				Expect(initial).To(Equal(expected))
				Expect(offset).To(Equal(expectedX))
			},
			Entry(
				"empty line",
				line{},
				5,
				line{},
				0,
			),
			Entry(
				"remove prefix fully",
				line{"styleA", "12345", "styleB", "6789"},
				7,
				line{"styleB", "89"},
				0,
			),
			Entry(
				"remove exact prefix",
				line{"styleA", "12345", "styleB", "6789"},
				5,
				line{"styleB", "6789"},
				0,
			),
			Entry(
				"remove exact prefix 2",
				line{"styleA", "12345", "styleB", "6789"},
				3,
				line{"styleA", "45", "styleB", "6789"},
				0,
			),
			Entry(
				"remove exact prefix 3",
				line{"styleA", "12345", "styleB", "6789"},
				4,
				line{"styleA", "5", "styleB", "6789"},
				0,
			),
			Entry(
				"remove more than line length",
				line{"styleA", "12345"},
				10,
				line{},
				0,
			),
		)
	})
	Context("CleanPostfix", func() {
		DescribeTable("Sunny",
			func(initial line, x int, expected line, expectedX int) {
				offset := initial.CleanPostfix(x)
				Expect(initial).To(Equal(expected))
				Expect(offset).To(Equal(expectedX))
			},
			Entry(
				"empty line",
				line{},
				5,
				line{},
				5,
			),
			Entry(
				"empty result",
				line{"styleA", "12345"},
				0,
				line{},
				0,
			),
			Entry(
				"one char prefix",
				line{"styleA", "12345"},
				1,
				line{"styleA", "1"},
				1,
			),
			Entry(
				"remove postfix fully",
				line{"styleA", "12345", "styleB", "6789"},
				7,
				line{"styleA", "12345", "styleB", "67"},
				7,
			),
			Entry(
				"remove exact postfix",
				line{"styleA", "12345", "styleB", "6789"},
				4,
				line{"styleA", "1234"},
				4,
			),
			Entry(
				"remove exact postfix",
				line{"styleA", "12345", "styleB", "6789"},
				5,
				line{"styleA", "12345"},
				5,
			),
			Entry(
				"remove exact postfix",
				line{"styleA", "12345", "styleB", "6789"},
				6,
				line{"styleA", "12345", "styleB", "6"},
				6,
			),
			Entry(
				"remove more than line length",
				line{"styleA", "12345"},
				10,
				line{"styleA", "12345"},
				10,
			),
		)
	})
})
