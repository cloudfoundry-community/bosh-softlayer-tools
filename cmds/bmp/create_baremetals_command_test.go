package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"

	fakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
)

var _ = Describe("create-baremetals command", func() {

	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command

		fakeBmpClient *fakes.FakeBmpClient
	)

	BeforeEach(func() {
		args = []string{"bmp", "create-baremetals"}
		options = cmds.Options{
			Verbose: false,
		}

		fakeBmpClient = fakes.NewFakeBmpClient("fake-username", "fake-password", "http://fake.url.com", "fake-config-path")
		cmd = bmp.NewCreateBaremetalsCommand(options, fakeBmpClient)
	})

	Describe("NewCreateBaremetalsCommand", func() {
		It("create new CreateBaremetalsCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewCreateBaremetalsCommand(options, fakeBmpClient)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a CreateBaremetalsCommand", func() {
			Expect(cmd.Name()).To(Equal("create-baremetals"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a CreateBaremetalsCommand", func() {
			Expect(cmd.Description()).To(Equal(`Create the missed baremetals: \"option --dryrun, only verify the orders\"`))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a CreateBaremetalsCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp create-baremetals [--dryrun]"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a CreateBaremetalsCommand", func() {
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
