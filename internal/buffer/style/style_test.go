package style

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Style", func() {
	var (
		s Style
	)
	BeforeEach(func() {
		s = newStyle()
	})
	Context("Style", func() {
		It("Sunny And Rainy", func() {
			Expect(s.String()).To(Equal(""))
			By("rgb bg")
			si := s.Set([]uint64{48, 2, 1, 10, 100})
			Expect(si.String()).To(Equal("\x1b[48;2;1;10;100m"))
			Expect(s.String()).To(Equal(""))
			By("rgb color")
			si = si.Set([]uint64{38, 2, 1, 10, 100})
			Expect(si.String()).To(Equal("\x1b[48;2;1;10;100m\x1b[38;2;1;10;100m"))
			By("bg")
			si = si.Set([]uint64{47})
			Expect(si.String()).To(Equal("\x1b[38;2;1;10;100m\x1b[47m"))
			By("style")
			si = si.Set([]uint64{36})
			Expect(si.String()).To(Equal("\x1b[47;36m"))
			By("color")
			si = si.Set([]uint64{37})
			Expect(si.String()).To(Equal("\x1b[47;37m"))
			By("style")
			si = si.Set([]uint64{2})
			Expect(si.String()).To(Equal("\x1b[47;37;2m"))
			si = si.Set([]uint64{0})
			Expect(si.String()).To(Equal(""))
		})
	})
})
