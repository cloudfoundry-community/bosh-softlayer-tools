package cmds_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
)

var _ = Describe("utils", func() {
	Describe("EqualOptions", func() {
		var (
			opts1, opts2 cmds.Options
		)

		BeforeEach(func() {
			opts1 = cmds.Options{
				Help:           false,
				Verbose:        true,
				DryRun:         false,
				Latest:         0,
				Packages:       "test-package",
				PackageOptions: "test-package-options",
			}

			opts2 = cmds.Options{
				Help:           false,
				Verbose:        true,
				DryRun:         false,
				Latest:         0,
				Packages:       "test-package",
				PackageOptions: "test-package-options",
			}
		})

		It("returns true for the same options", func() {
			Expect(cmds.EqualOptions(opts1, opts1)).To(BeTrue())
			Expect(cmds.EqualOptions(opts2, opts2)).To(BeTrue())
		})

		It("returns true for two equal options", func() {
			Expect(cmds.EqualOptions(opts1, opts2)).To(BeTrue())
		})

		It("returns false when options2 is modified", func() {
			opts2.Verbose = false
			Expect(cmds.EqualOptions(opts1, opts2)).To(BeFalse())
		})

		It("returns false when one options is empty options", func() {
			Expect(cmds.EqualOptions(cmds.Options{}, opts2)).To(BeFalse())
			Expect(cmds.EqualOptions(opts1, cmds.Options{})).To(BeFalse())
		})
	})
})
