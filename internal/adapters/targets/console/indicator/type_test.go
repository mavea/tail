package indicator

import (
	"tail/internal/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Type", func() {
	var (
		t config.StringValue
	)
	BeforeEach(func() {
		t = NewType()
	})
	Context("Type", func() {
		It("Sunny And Rainy", func() {
			Expect(t).To(Not(BeNil()))
			Expect(t.String()).To(Equal("none"))
			Expect(t.Set("roll")).To(BeNil())
			Expect(t.String()).To(Equal("roll"))
			err := t.Set("roll1")
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(Equal("invalid screen type value: roll1"))
			Expect(t.String()).To(Equal("roll"))
		})
	})

})
