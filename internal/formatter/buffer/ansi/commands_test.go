package ansi

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Action", func() {
	var (
		tDECPrivateModes   DECPrivateModes
		tSecondaryDA       SecondaryDA
		tWindowAndPalettes WindowAndPalettes
		tMainList          MainList
		tOtherList         OtherList
	)
	BeforeEach(func() {
		tDECPrivateModes = GetActionsDECPrivateModes()
		tSecondaryDA = GetActionsSecondaryDA()
		tWindowAndPalettes = GetActionsWindowAndPalettes()
		tMainList = GetActionsMainList()
		tOtherList = GetActionsOtherList()
	})
	Context("DECPrivateModes", func() {
		It("Sunny And Rainy", func() {
			Expect(tDECPrivateModes).To(Not(BeNil()))
			f, ok := tDECPrivateModes.GetFunction('h')
			Expect(ok).To(Equal(true))
			Expect(f).To(Not(BeNil()))
			com, ok2 := f([]uint64{1, 2, 3})
			Expect(com).To(BeNil())
			Expect(ok2).To(Equal(false))
			com, ok2 = f([]uint64{1})
			Expect(ok2).To(Equal(true))
			Expect(com.GetType()).To(Equal(KeyOtherDo))
			Expect(com.GetValsUint()).To(Equal([]uint64{1}))
			Expect(com.GetValsString()).To(Equal(""))
		})
	})
	Context("SecondaryDA", func() {
		It("Sunny And Rainy", func() {
			Expect(tSecondaryDA).To(Not(BeNil()))
			f, ok := tSecondaryDA.GetFunction('c')
			Expect(ok).To(Equal(true))
			Expect(f).To(Not(BeNil()))
			com, ok2 := f([]uint64{1, 2, 3})
			Expect(ok2).To(Equal(false))
			Expect(com).To(BeNil())
			com, ok2 = f([]uint64{1})
			Expect(ok2).To(Equal(true))
			Expect(com.GetType()).To(Equal(KeyTerminalID))
			Expect(com.GetValsUint()).To(Equal([]uint64{1, 0}))
			Expect(com.GetValsString()).To(Equal(""))
		})
	})
	Context("WindowAndPalettes", func() {
		It("Sunny And Rainy", func() {
			Expect(tWindowAndPalettes).To(Not(BeNil()))
			f, ok := tWindowAndPalettes.GetFunction('s')
			Expect(ok).To(Equal(true))
			Expect(f).To(Not(BeNil()))
			com, ok2 := f([]uint64{1, 2, 3}, "test")
			Expect(ok2).To(Equal(false))
			Expect(com).To(BeNil())
			com, ok2 = f([]uint64{1, 7}, "test")
			Expect(ok2).To(Equal(true))
			Expect(com.GetType()).To(Equal(KeySetPalettes))
			Expect(com.GetValsUint()).To(Equal([]uint64{1, 7}))
			Expect(com.GetValsString()).To(Equal(""))
		})
	})
	Context("MainList", func() {
		It("Sunny And Rainy", func() {
			Expect(tMainList).To(Not(BeNil()))
			f, ok := tMainList.GetFunction('D')
			Expect(ok).To(Equal(true))
			Expect(f).To(Not(BeNil()))
			com, ok2 := f([]uint64{1})
			Expect(ok2).To(Equal(true))
			Expect(com.GetType()).To(Equal(KeyCursorLeft))
			Expect(com.GetValsUint()).To(Equal([]uint64{1}))
			Expect(com.GetValsString()).To(Equal(""))
		})
	})
	Context("OtherList", func() {
		It("Sunny And Rainy", func() {
			Expect(tOtherList).To(Not(BeNil()))
			Expect(len(tOtherList)).To(Equal(0))
		})
	})

})
