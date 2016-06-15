package bmp_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("provisioning-baremetal command", func() {
	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command

		fakeBmpClient *fakes.FakeBmpClient

		userInput string
	)

	BeforeEach(func() {
		args = []string{"bmp", "provisiong-baremetal"}
		options = cmds.Options{
			Verbose:      false,
			Stemecell:    "fake-stemcell",
			VMPrefix:     "fake-vmprefix",
			NetbootImage: "fake-netboot-image",
		}

		fakeBmpClient = fakes.NewFakeBmpClient("fake-username", "fake-password", "http://fake.url.com", "fake-config-path")
		cmd = bmp.NewProvisioningBaremetalCommand(options, fakeBmpClient)
	})

	Describe("NewProvisioningBaremetalCommand", func() {
		It("create new ProvisioningBaremetalCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewProvisioningBaremetalCommand(options, fakeBmpClient)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a ProvisioningBaremetalCommand", func() {
			Expect(cmd.Name()).To(Equal("provisioning-baremetal"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a ProvisioningBaremetalCommand", func() {
			Expect(cmd.Description()).To(Equal(`provisioning a baremetal with specific stemcell, netboot image`))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a ProvisioningBaremetalCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp provisioning-baremetal --vmprefix <vm-prefix> --stemcell <bm-stemcell> --netbootimage <bm-netboot-image>"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a ProvisioningBaremetalCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
			Expect(cmd.Options().VMPrefix).To(Equal("fake-vmprefix"))
			Expect(cmd.Options().Stemecell).To(Equal("fake-stemcell"))
			Expect(cmd.Options().NetbootImage).To(Equal("fake-netboot-image"))
		})
	})

	Describe("#Validate", func() {
		It("validates a good ProvisioningBaremetalCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when validating a bad ProvisioningBaremetalCommand", func() {
			Context("when no options passed", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose: false,
					}
				})

				It("fails validation with errors", func() {
					cmd = bmp.NewProvisioningBaremetalCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when no stemcell passed", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:      false,
						VMPrefix:     "fake-vmprefix",
						NetbootImage: "fake-netboot-image",
					}
				})

				It("fails validation with errors", func() {
					cmd = bmp.NewProvisioningBaremetalCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when no vmprefix passed", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:      false,
						Stemecell:    "fake-stemcell",
						NetbootImage: "fake-netboot-image",
					}
				})

				It("fails validation with errors", func() {
					cmd = bmp.NewProvisioningBaremetalCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when no netboot image passed", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:   false,
						VMPrefix:  "fake-vmprefix",
						Stemecell: "fake-stemcell",
					}
				})

				It("fails validation with errors", func() {
					cmd = bmp.NewProvisioningBaremetalCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

		})
	})

	Describe("#Execute", func() {
		Context("when executing a good ProvisioningBaremetalCommand", func() {
			BeforeEach(func() {
				fakeBmpClient.ProvisioningBaremetalResponse.Status = 200
				fakeBmpClient.ProvisioningBaremetalErr = nil
				userInput = "yes"
				cmd = bmp.NewFakeProvisioningBaremetalCommand(options, fakeBmpClient, userInput)
			})

			It("executes a good ProvisioningBaremetalCommand", func() {
				rc, err := cmd.Execute(args)
				Expect(rc).To(Equal(0))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when executing a bad ProvisioningBaremetalCommand", func() {
			Context("when confirmation fails", func() {
				BeforeEach(func() {
					userInput = "no"
					cmd = bmp.NewFakeProvisioningBaremetalCommand(options, fakeBmpClient, userInput)
				})

				It("executes with error", func() {
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(1))
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when executing ProvisioningBaremetalCommand with error", func() {
				BeforeEach(func() {
					fakeBmpClient.ProvisioningBaremetalResponse.Status = 500
					fakeBmpClient.ProvisioningBaremetalErr = errors.New("500")
					userInput = "yes"
					cmd = bmp.NewFakeProvisioningBaremetalCommand(options, fakeBmpClient, userInput)
				})

				It("executes with error", func() {
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(1))
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when ProvisioningBaremetalCommand response different than 200", func() {
				BeforeEach(func() {
					fakeBmpClient.ProvisioningBaremetalResponse.Status = 404
					userInput = "yes"
					cmd = bmp.NewFakeProvisioningBaremetalCommand(options, fakeBmpClient, userInput)
				})

				It("executes with error", func() {
					rc, err := cmd.Execute(args)
					Expect(rc).To(Equal(404))
					Expect(err).ToNot(HaveOccurred())
				})
			})

		})

	})
})
