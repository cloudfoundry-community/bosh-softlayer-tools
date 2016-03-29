package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("sl command", func() {

	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command
	)

	BeforeEach(func() {
		args = []string{"bmp", "sl"}
		options = cmds.Options{
			Verbose: false,
		}

		cmd = bmp.NewSlCommand(options)
	})

	Describe("NewSlCommand", func() {
		It("create new SlCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewSlCommand(options)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a SlCommand", func() {
			Expect(cmd.Name()).To(Equal("sl"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a SlCommand", func() {
			Expect(cmd.Description()).To(Equal("List all Softlayer packages or package options"))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a SlCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp sl --packages | --package-options"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a SlCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
		})
	})

	Describe("#Validate", func() {
		It("validates a good SlCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("#Execute", func() {
		It("executes a good SlCommand", func() {
			rc, err := cmd.Execute(args)
			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
