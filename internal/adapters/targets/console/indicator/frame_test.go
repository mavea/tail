package indicator

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Icon", func() {
	var (
		i Indicator
	)
	BeforeEach(func() {
		i = New("")
	})
	Context("Icon", func() {
		It("Sunny And Rainy", func() {
			Expect(i).To(Not(BeNil()))
			Expect(i.Get()).To(Equal(""))
			Expect(i.Get()).To(Equal(""))
			Expect(i.Clean()).To(Equal(""))
			i = New("roll")
			Expect(i).To(Not(BeNil()))
			g := i.Get()
			Expect(g).To(Not(Equal("")))
			Expect(i.Get()).To(Not(Equal(g)))
			Expect(i.Clean()).To(Equal(""))
		})
	})
})
