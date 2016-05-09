package bmp_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var _ = Describe("sl command", func() {

	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command

		fakeBmpClient *fakes.FakeBmpClient
	)

	BeforeEach(func() {
		args = []string{"bmp", "sl"}
		options = cmds.Options{
			Verbose: false,
		}

		fakeBmpClient = fakes.NewFakeBmpClient("fake-username", "fake-password", "http://fake.url.com", "fake-config-path")
		cmd = bmp.NewSlCommand(options, fakeBmpClient)
	})

	Describe("NewSlCommand", func() {
		It("create new SlCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewSlCommand(options, fakeBmpClient)
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
		Context("validates a good TaskCommand", func() {
			Context("when --pacakges specified", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:  false,
						Packages: true,
					}
				})

				It("validates without errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeTrue())
					Expect(err).ToNot(HaveOccurred())
				})
			})
			//TODO: Add test cases for --package-options
		})

		Context("validates a bad SlCommand", func() {
			Context("when neither --packages nor --package-options specified", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose: false,
					}
				})

				It("validates with errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

			//TODO: Add test cases for --package-options

		})
	})

	Describe("#Execute", func() {
		Context("executes a good SlCommand", func() {
			Context("executes sl --packages", func() {
				BeforeEach(func() {
					fakeBmpClient.SlPackagesResponse.Status = 200
					fakeBmpClient.SlPackagesErr = nil
					options = cmds.Options{
						Verbose:  false,
						Packages: true,
					}
				})
				It("executes without errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					rc, err := cmd.Execute([]string{"bmp", "sl", "--packages"})
					Expect(rc).To(Equal(0))
					Expect(err).ToNot(HaveOccurred())
				})
			})

			//TODO: Add test cases for --package-options
		})

		Context("executes a bad SlCommand", func() {
			Context("when SlCommand --packages fails", func() {
				BeforeEach(func() {
					fakeBmpClient.SlPackagesResponse.Status = 500
					fakeBmpClient.SlPackagesErr = errors.New("500")
					options = cmds.Options{
						Verbose:  false,
						Packages: true,
					}
				})

				It("executes with error", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					rc, err := cmd.Execute([]string{"bmp", "sl", "--packages"})
					Expect(rc).To(Equal(1))
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when SlCommand --packages response different than 200", func() {
				BeforeEach(func() {
					fakeBmpClient.SlPackagesResponse.Status = 404
					options = cmds.Options{
						Verbose:  false,
						Packages: true,
					}
				})

				It("executes without error", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					rc, err := cmd.Execute([]string{"bmp", "sl", "--packages"})
					Expect(rc).To(Equal(404))
					Expect(err).ToNot(HaveOccurred())
				})
			})

			//TODO: Add test cases for --package-options
		})
	})
})
