package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("bms command", func() {

	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command
	)

	BeforeEach(func() {
		args = []string{"bmp", "bms"}
		options = cmds.Options{
			Verbose: false,
		}

		cmd = bmp.NewBmsCommand(options)
	})

	Describe("NewBmsCommand", func() {
		It("create new BmsCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewBmsCommand(options)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a BmsCommand", func() {
			Expect(cmd.Name()).To(Equal("bms"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a BmsCommand", func() {
			Expect(cmd.Description()).To(Equal("List all bare metals"))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a BmsCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp bms"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a BmsCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
		})
	})

	Describe("#Validate", func() {
		It("validates a good BmsCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("#Execute", func() {
		It("executes a good BmsCommand", func() {
			rc, err := cmd.Execute(args)
			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
