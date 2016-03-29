package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("target command", func() {
	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command
	)

	BeforeEach(func() {
		args = []string{"bmp", "target"}
		options = cmds.Options{
			Verbose: false,
		}

		cmd = bmp.NewTargetCommand(options)
	})

	Describe("NewTargetCommand", func() {
		It("create new TargetCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewTargetCommand(options)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a TargetCommand", func() {
			Expect(cmd.Name()).To(Equal("target"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a TargetCommand", func() {
			Expect(cmd.Description()).To(Equal("Set the URL of Bare Metal Provision Server"))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a TargetCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp target http://url.to.bmp.server"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a TargetCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
		})
	})

	Describe("#Validate", func() {
		It("validates a good TargetCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("#Execute", func() {
		It("executes a good TargetCommand", func() {
			rc, err := cmd.Execute(args)
			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
