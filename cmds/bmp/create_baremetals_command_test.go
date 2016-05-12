package bmp_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
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
			Verbose:    false,
			Deployment: "../../test_fixtures/bmp/deployment.yml",
			DryRun:     false,
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
			Expect(cmd.Usage()).To(Equal("bmp create-baremetals --deployment[-d] <deployment file> [--dryrun]"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a CreateBaremetalsCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
			Expect(cmd.Options().Deployment).ToNot(Equal(""))
			Expect(cmd.Options().Deployment).To(Equal("../../test_fixtures/bmp/deployment.yml"))
		})
	})

	Describe("#Validate", func() {
		It("validates a good CreateBaremetalsCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		Context("validates a bad CreateBaremetalsCommand", func() {
			BeforeEach(func() {
				options = cmds.Options{
					Verbose:    false,
					Deployment: "fake-deployment-file",
				}
			})

			It("fails validation when deployment file not existed", func() {
				cmd = bmp.NewCreateBaremetalsCommand(options, fakeBmpClient)
				validate, err := cmd.Validate()
				Expect(validate).To(BeFalse())
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("#Execute", func() {
		Context("executes a good CreateBaremetalsCommand", func() {
			Context("when executes CreateBaremetalsCommand without --dryrun", func() {
				BeforeEach(func() {
					fakeBmpClient.CreateBaremetalsResponse.Status = 200
					fakeBmpClient.CreateBaremetalsErr = nil
				})

				It("executes with no error", func() {
					args = []string{"bmp", "create-baremetals", "-d", "../../test_fixtures/bmp/deployment.yml"}
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(0))
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("when executes CreateBaremetalsCommand with --dryrun", func() {
				BeforeEach(func() {
					fakeBmpClient.CreateBaremetalsResponse.Status = 200
					fakeBmpClient.CreateBaremetalsErr = nil
					options = cmds.Options{
						Verbose:    false,
						Deployment: "../../test_fixtures/bmp/deployment.yml",
						DryRun:     true,
					}
				})

				It("executes with no error", func() {
					args = []string{"bmp", "create-baremetals", "-d", "../../test_fixtures/bmp/deployment.yml", "--dryrun"}
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(0))
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Context("executes a bad CreateBaremetalsCommand", func() {
			Context("when CreateBaremetalsCommand fails", func() {
				BeforeEach(func() {
					fakeBmpClient.CreateBaremetalsResponse.Status = 500
					fakeBmpClient.CreateBaremetalsErr = errors.New("500")
				})

				It("executes with errors", func() {
					args = []string{"bmp", "create-baremetals", "-d", "../../test_fixtures/bmp/deployment.yml"}
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(1))
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when CreateBaremetalsCommand response different than 200", func() {
				BeforeEach(func() {
					fakeBmpClient.CreateBaremetalsResponse.Status = 404
				})

				It("executes without errors", func() {
					args = []string{"bmp", "create-baremetals", "-d", "../../test_fixtures/bmp/deployment.yml"}
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(404))
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
