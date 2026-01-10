package line

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lines", func() {
	Context("Add", func() {
		DescribeTable("Sunny",
			func(initial lines, count uint64, expected lines) {
				initial.Add(count)
				Expect(initial).To(Equal(expected))
			},
			Entry("add zero", lines{}, uint64(0), lines{}),
			Entry("add one", lines{}, uint64(1), lines{NewLine()}),
			Entry("add multiple", lines{}, uint64(3), lines{NewLine(), NewLine(), NewLine()}),
			Entry("add to existing", lines{NewLine()}, uint64(2), lines{NewLine(), NewLine(), NewLine()}),
		)
	})
	Context("Get", func() {
		var list = &lines{&line{
			"styleA", "contentA",
			"styleB", "contentB",
			"styleC", "contentC",
		}, &line{
			"styleAA", "contentAA",
			"styleBA", "contentBA",
			"styleCA", "contentCA",
		}, &line{
			"styleAB", "contentAB",
			"styleBB", "contentBB",
			"styleCB", "contentCB",
		}}
		DescribeTable("Sunny",
			func(initial *lines, id uint64, expected *line) {
				Expect(initial.Get(id)).To(Equal(expected))
			},
			Entry("get first",
				list,
				uint64(0),
				&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}),
			Entry("get middle",
				list,
				uint64(1),
				&line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}),
			Entry("get last",
				list,
				uint64(2),
				&line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}),
		)
	})
	Context("CleanPostfix", func() {
		DescribeTable("Sunny",
			func(count uint64, expected *lines) {
				var list = &lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}, &line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}, &line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}
				list.CleanPostfix(count)
				Expect(list).To(Equal(expected))
			},
			Entry("clean nothing",
				uint64(2),
				&lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}, &line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}}),
			Entry("clean one",
				uint64(1),
				&lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}}),
			Entry("clear all",
				uint64(0),
				&lines{}),
			Entry("no clean",
				uint64(10),
				&lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}, &line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}, &line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}),
		)
	})
	Context("CleanPrefix", func() {
		DescribeTable("Sunny",
			func(count uint64, expected *lines) {
				var list = &lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}, &line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}, &line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}
				list.CleanPrefix(count)
				Expect(list).To(Equal(expected))
			},
			Entry("clean two",
				uint64(2),
				&lines{&line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}),
			Entry("clean one",
				uint64(1),
				&lines{&line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}, &line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}),
			Entry("clean nothing",
				uint64(3),
				&lines{}),
			Entry("no clean",
				uint64(0),
				&lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}, &line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}, &line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}),
		)
	})
	Context("CleanString", func() {
		DescribeTable("Sunny",
			func(count uint64, expected *lines) {
				var list = &lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}, &line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}, &line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}
				list.CleanString(count)
				Expect(list).To(Equal(expected))
			},
			Entry("clean first string",
				uint64(0),
				&lines{&line{}, &line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}, &line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}),
			Entry("clean middle string",
				uint64(1),
				&lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}, &line{}, &line{
					"styleAB", "contentAB",
					"styleBB", "contentBB",
					"styleCB", "contentCB",
				}}),
			Entry("clean last string",
				uint64(2),
				&lines{&line{
					"styleA", "contentA",
					"styleB", "contentB",
					"styleC", "contentC",
				}, &line{
					"styleAA", "contentAA",
					"styleBA", "contentBA",
					"styleCA", "contentCA",
				}, &line{}}),
		)
	})
	Context("Len", func() {
		DescribeTable("Sunny",
			func(initial lines, expected uint64) {
				Expect(initial.Len()).To(Equal(expected))
			},
			Entry("empty lines", lines{}, uint64(0)),
			Entry("single line", lines{NewLine()}, uint64(1)),
			Entry("multiple lines", lines{NewLine(), NewLine()}, uint64(2)),
		)
	})
	Context("NewLines", func() {
		DescribeTable("Sunny",
			func(initial uint64, expected *lines) {
				Expect(NewLines(initial)).To(Equal(expected))
			},
			Entry("zero lines", uint64(0), &lines{}),
			Entry("one lines", uint64(1), &lines{NewLine()}),
			Entry("multiple lines", uint64(3), &lines{NewLine(), NewLine(), NewLine()}),
		)
	})
})
