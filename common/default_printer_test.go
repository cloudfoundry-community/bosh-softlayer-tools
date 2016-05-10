package common_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
	fakes "github.com/cloudfoundry-community/bosh-softlayer-tools/common/fakes"
)

var _ = Describe("DefaultPrinter", func() {
	var (
		printer common.Printer
		fakeUi  *fakes.FakeUi
	)

	BeforeEach(func() {
		fakeUi = fakes.NewFakeUi()
	})

	Describe("when verbose is true", func() {
		BeforeEach(func() {
			printer = common.NewDefaultPrinter(fakeUi, true)
		})

		It("#Printf", func() {
			rc, err := printer.Printf("%s %s", "hello", "world")

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeUi.Output).To(ContainSubstring(fmt.Sprintf("%s %s", "hello", "world")))
		})

		It("#Println", func() {
			rc, err := printer.Println("hello")

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeUi.Output).To(ContainSubstring(fmt.Sprint("hello")))
		})

		It("#PrintfInfo", func() {
			rc, err := printer.PrintfInfo("%s %s", "hello", "world")

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeUi.Output).To(ContainSubstring(fmt.Sprint("hello world")))
		})

		It("#PrintlnInfo", func() {
			rc, err := printer.PrintlnInfo("hello")

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeUi.Output).To(ContainSubstring(fmt.Sprint("hello")))
		})
	})

	Describe("when verbose is false", func() {
		BeforeEach(func() {
			printer = common.NewDefaultPrinter(fakeUi, false)
		})

		It("#Printf", func() {
			rc, err := printer.Printf("%s %s", "hello", "world")

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeUi.Output).To(Equal(""))
		})

		It("#Println", func() {
			rc, err := printer.Println("hello")

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeUi.Output).To(Equal(""))
		})

		It("#PrintfInfo", func() {
			rc, err := printer.PrintfInfo("%s %s", "hello", "world")

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeUi.Output).To(ContainSubstring(fmt.Sprint("hello world")))
		})

		It("#PrintlnInfo", func() {
			rc, err := printer.PrintlnInfo("hello")

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeUi.Output).To(ContainSubstring(fmt.Sprint("hello")))
		})
	})
})
