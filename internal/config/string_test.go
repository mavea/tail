package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultString", func() {
	Context("String", func() {
		DescribeTable("Sunny",
			func(initial, expected string) {
				ds := &defaultString{value: initial}
				Expect(ds.String()).To(Equal(expected))
			},
			Entry("default string", "hello", "hello"),
			Entry("empty string", "", ""),
		)
	})
	Context("Set", func() {
		DescribeTable("Sunny",
			func(initial string) {
				ds := &defaultString{}
				err := ds.Set(initial)
				Expect(err).To(BeNil())
			},
			Entry("valid string", "newValue"),
			Entry("empty string", ""),
		)
	})
	Context("Type", func() {
		ds := &defaultString{}
		expected := "string"
		Expect(ds.Type()).To(Equal(expected))
	})
})
