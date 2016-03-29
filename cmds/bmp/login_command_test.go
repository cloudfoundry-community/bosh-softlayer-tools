package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("login command", func() {

	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command
	)

	BeforeEach(func() {
		args = []string{"bmp", "login"}
		options = cmds.Options{
			Verbose: false,
		}

		cmd = bmp.NewLoginCommand(options)
	})

	Describe("NewLoginCommand", func() {
		It("create new LoginCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewLoginCommand(options)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a LoginCommand", func() {
			Expect(cmd.Name()).To(Equal("login"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a LoginCommand", func() {
			Expect(cmd.Description()).To(Equal("Login to the Bare Metal Provision Server"))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a LoginCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp login"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a LoginCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
		})
	})

	Describe("#Validate", func() {
		It("validates a good LoginCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("#Execute", func() {
		It("executes a good LoginCommand", func() {
			rc, err := cmd.Execute(args)
			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
