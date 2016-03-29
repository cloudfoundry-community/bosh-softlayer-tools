package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("stemcells command", func() {
	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command
	)

	BeforeEach(func() {
		args = []string{"bmp", "stemcells"}
		options = cmds.Options{
			Verbose: false,
		}

		cmd = bmp.NewStemcellsCommand(options)
	})

	Describe("NewStemcellsCommand", func() {
		It("create new StemcellsCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewStemcellsCommand(options)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a StemcellsCommand", func() {
			Expect(cmd.Name()).To(Equal("stemcells"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a StemcellsCommand", func() {
			Expect(cmd.Description()).To(Equal("List all stemcells"))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a StemcellsCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp stemcells"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a StemcellsCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
		})
	})

	Describe("#Validate", func() {
		It("validates a good StemcellsCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("#Execute", func() {
		It("executes a good StemcellsCommand", func() {
			rc, err := cmd.Execute(args)
			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
