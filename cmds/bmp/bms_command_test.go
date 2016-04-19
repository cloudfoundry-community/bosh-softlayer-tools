package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"

	fakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
	"errors"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
)

var _ = Describe("bms command", func() {

	var (
		args    []string
		options cmds.Options
		cmd     cmds.Command

		fakeBmpClient *fakes.FakeBmpClient
	)

	BeforeEach(func() {
		args = []string{"bmp", "bms"}
		options = cmds.Options{
			Verbose: false,
			Deployment: "../../test_fixtures/bmp/deployment.yml",
		}

		fakeBmpClient = fakes.NewFakeBmpClient("fake-username", "fake-password", "http://fake.url.com", "fake-config-path")
		cmd = bmp.NewBmsCommand(options, fakeBmpClient)
	})

	Describe("NewBmsCommand", func() {
		It("create new BmsCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewBmsCommand(options, fakeBmpClient)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a BmsCommand", func() {
			Expect(cmd.Name()).To(Equal("bms"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a BmsCommand", func() {
			Expect(cmd.Description()).To(Equal("List all bare metals"))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a BmsCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp bms --deployment[-d] <deployment file>"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a BmsCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())

			Expect(cmd.Options().Deployment).ToNot(Equal(""))
			Expect(cmd.Options().Deployment).To(Equal("../../test_fixtures/bmp/deployment.yml"))
		})
	})

	Describe("#Validate", func() {
		It("validates a good BmsCommand", func() {
			validate, err := cmd.Validate()
			Expect(validate).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})

		Context("bad BmsCommand", func() {
			Context("no deployment file", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose: false,
						Deployment: "",
					}
				})

				It("fails validation", func(){
					cmd = bmp.NewBmsCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("deployment file doesn't exist", func() {
				BeforeEach(func() {
					options = cmds.Options{
						Verbose: false,
						Deployment: "fake-deployment-file",
					}
				})

				It("fails validation", func(){
					cmd = bmp.NewBmsCommand(options, fakeBmpClient)
					validate, err := cmd.Validate()
					Expect(validate).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("#Execute", func() {
		Context("executes a good BmsCommand", func() {
			BeforeEach(func() {
				fakeBmpClient.BmsResponse.Status = 200
				fakeBmpClient.BmsErr = nil
				fakeBmpClient.BmsResponse.Data = []clients.BaremetalInfo{
					clients.BaremetalInfo{
						Id:   0,
						Hostname: "hostname0",
						Private_ip_address: "private_ip_address0",
						Public_ip_address: "public_ip_address0",
						Hardware_status: "hardware_status0",
						Memory: 0,
						Cpu: 0,
						Provision_date: "provision_date0",
					}}
			})

			It("executes with no error", func() {
				rc, err := cmd.Execute(args)
				Expect(rc).To(Equal(0))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("executes a bad Bmscommand", func() {
			BeforeEach(func() {
				fakeBmpClient.BmsResponse.Status = 404
				fakeBmpClient.BmsErr = errors.New("fake-error")
			})

			It("executes with error", func() {
				rc, err := cmd.Execute(args)
				Expect(rc).To(Equal(404))
				Expect(err).To(HaveOccurred())
			})

		})

	})
})
