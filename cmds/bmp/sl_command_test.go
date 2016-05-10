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
		Context("validates a good SlCommand", func() {
			Context("when --packages specified", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:  false,
						Packages: true,
					}
				})

				It("validation pass without errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeTrue())
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("when --package-options specified", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:        false,
						PackageOptions: "1",
					}
				})

				It("validation pass without errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeTrue())
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Context("validates a bad SlCommand", func() {
			Context("when neither --packages nor --package-options specified", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose: false,
					}
				})

				It("validation fails with errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when --package-options specified but no value", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:        false,
						PackageOptions: "",
					}
				})

				It("validation fails with errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("#Execute", func() {
		Context("executes a good SlCommand", func() {
			Context("when executes sl --packages", func() {
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

			Context("when executes sl --package-options 1", func() {
				BeforeEach(func() {
					fakeBmpClient.SlPackageOptionsResponse.Status = 200
					fakeBmpClient.SlPackageOptionsErr = nil
					options = cmds.Options{
						Verbose:        false,
						PackageOptions: "1",
					}
				})

				It("executes without errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					rc, err := cmd.Execute([]string{"bmp", "sl", "--package-options", "1"})
					Expect(rc).To(Equal(0))
					Expect(err).ToNot(HaveOccurred())
				})
			})
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

				It("executes with errors", func() {
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

				It("executes without errors", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					rc, err := cmd.Execute([]string{"bmp", "sl", "--packages"})
					Expect(rc).To(Equal(404))
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("when SlCommand --package-options fails", func() {
				BeforeEach(func() {
					fakeBmpClient.SlPackageOptionsResponse.Status = 500
					fakeBmpClient.SlPackageOptionsErr = errors.New("500")
					options = cmds.Options{
						Verbose:        false,
						PackageOptions: "1",
					}
				})

				It("executes with error", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					rc, err := cmd.Execute([]string{"bmp", "sl", "--package-options", "1"})
					Expect(rc).To(Equal(1))
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when SlCommand --package-options response different than 200", func() {
				BeforeEach(func() {
					fakeBmpClient.SlPackageOptionsResponse.Status = 404
					options = cmds.Options{
						Verbose:        false,
						PackageOptions: "1",
					}
				})

				It("executes without error", func() {
					cmd = bmp.NewSlCommand(options, fakeBmpClient)
					rc, err := cmd.Execute([]string{"bmp", "sl", "--package-options", "1"})
					Expect(rc).To(Equal(404))
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
