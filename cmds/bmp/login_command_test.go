package bmp_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"

	clientsfakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
)

var _ = Describe("login command", func() {

	var (
		args          []string
		options       cmds.Options
		cmd           cmds.Command
		fakeBmpClient *clientsfakes.FakeBmpClient
	)

	BeforeEach(func() {
		args = []string{"bmp", "login"}
		options = cmds.Options{
			Verbose:  false,
			Username: "fake-username",
			Password: "fake-password",
		}

		fakeBmpClient = clientsfakes.NewFakeBmpClient(options.Username, options.Password, "http://fake.target.url")
		cmd = bmp.NewLoginCommand(options, fakeBmpClient)
	})

	Describe("NewLoginCommand", func() {
		It("create new LoginCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewLoginCommand(options, fakeBmpClient)
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
			Expect(cmd.Usage()).To(Equal("bmp login --username[-u] <username> --password[-p] <password"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a LoginCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())

			Expect(cmd.Options().Username).ToNot(Equal(""))
			Expect(cmd.Options().Username).To(Equal("fake-username"))

			Expect(cmd.Options().Password).ToNot(Equal(""))
			Expect(cmd.Options().Password).To(Equal("fake-password"))
		})
	})

	Describe("#Validate", func() {
		It("validates a good LoginCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		Context("bad LoginCommand", func() {
			Context("no Username", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:  false,
						Username: "",
						Password: "fake-password",
					}
				})

				It("fails validation", func() {
					cmd = bmp.NewLoginCommand(options, fakeBmpClient)

					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("no Password", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:  false,
						Username: "fake-username",
						Password: "",
					}
				})

				It("fails validation", func() {
					cmd = bmp.NewLoginCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("no Username and no Password", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose:  false,
						Username: "",
						Password: "",
					}
				})

				It("fails validation", func() {
					cmd = bmp.NewLoginCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("#Execute", func() {
		Context("good LoginCommand", func() {
			BeforeEach(func() {
				fakeBmpClient.LoginResponse.Status = 200
				fakeBmpClient.LoginErr = nil
			})

			It("executes with no error", func() {
				rc, err := cmd.Execute(args)
				Expect(rc).To(Equal(0))
				Expect(err).ToNot(HaveOccurred())
			})

			It("saves the Username and Password to Config", func() {
				//TODO: @grace
			})
		})

		Context("bad LoginCommand", func() {
			BeforeEach(func() {
				fakeBmpClient.LoginResponse.Status = 500
				fakeBmpClient.LoginErr = errors.New("500")
			})

			It("executes with error", func() {
				rc, err := cmd.Execute(args)
				Expect(rc).To(Equal(500))
				Expect(err).To(HaveOccurred())
			})

			//TODO: verify LoginResponse.Status different than 200
			//TODO: verify Login() execution failing
		})
	})
})
