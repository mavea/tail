package indicator

import (
	configGeneral "tail/internal/config/general"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Type", func() {
	var (
		t configGeneral.StringValue
	)
	BeforeEach(func() {
		t = NewIndicatorType()
	})
	Context("Type", func() {
		It("Sunny And Rainy", func() {
			Expect(t).To(Not(BeNil()))
			Expect(t.Type()).To(Equal("string"))
			Expect(t.Validate()).To(Equal(true))
			Expect(t.String()).To(Equal("none"))
			Expect(t.Set("roller")).To(BeNil())
			Expect(t.String()).To(Equal("roller"))
			Expect(t.Validate()).To(Equal(true))
			err := t.Set("roll1")
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(Equal("invalid screen type value: roll1"))
			Expect(t.String()).To(Equal("roller"))
		})
	})

})
