package bmp_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("stemcells command", func() {
	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command

		fakeBmpClient *fakes.FakeBmpClient
	)

	BeforeEach(func() {
		args = []string{"bmp", "stemcells"}
		options = cmds.Options{
			Verbose: false,
		}

		fakeBmpClient = fakes.NewFakeBmpClient("fake-username", "fake-password", "http://fake.url.com", "fake-config-path")
		cmd = bmp.NewStemcellsCommand(options, fakeBmpClient)
	})

	Describe("NewStemcellsCommand", func() {
		It("create new StemcellsCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewStemcellsCommand(options, fakeBmpClient)
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
		Context("executes a good StemcellsCommand", func() {
			BeforeEach(func() {
				fakeBmpClient.StemcellResponse.Status = 200
				fakeBmpClient.StemcellErr = nil
			})

			It("execute with no error", func() {
				rc, err := cmd.Execute(args)
				Expect(rc).To(Equal(0))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("executes a bad StemcellsCommand", func() {
			Context("executes StemcellsCommand with error", func() {
				BeforeEach(func() {
					fakeBmpClient.StemcellResponse.Status = 500
					fakeBmpClient.StemcellErr = errors.New("500")
				})

				It("executes with error", func() {
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(1))
					Expect(err).To(HaveOccurred())
				})
			})

			Context("StemcellsCommand response different than 200", func() {
				BeforeEach(func() {
					fakeBmpClient.StemcellResponse.Status = 404
				})

				It("response code different than 200", func() {
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(404))
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
