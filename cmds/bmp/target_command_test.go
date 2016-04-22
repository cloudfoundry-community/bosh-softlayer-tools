package bmp_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
	config "github.com/cloudfoundry-community/bosh-softlayer-tools/config"

	clientsfakes "github.com/cloudfoundry-community/bosh-softlayer-tools/clients/fakes"
)

var _ = Describe("target command", func() {
	var (
		err error

		args    []string
		options cmds.Options
		cmd     cmds.Command

		tmpDir, tmpFileName string

		fakeBmpClient *clientsfakes.FakeBmpClient
	)

	BeforeEach(func() {
		args = []string{"bmp", "target"}
		options = cmds.Options{
			Verbose: false,
			Target:  "http://fake.url",
		}

		tmpDir, err = ioutil.TempDir("", "bmp-target-execute")
		Expect(err).ToNot(HaveOccurred())

		tmpFile, err := ioutil.TempFile(tmpDir, ".bmp_config")
		Expect(err).ToNot(HaveOccurred())

		tmpFileName = tmpFile.Name()

		fakeBmpClient = clientsfakes.NewFakeBmpClient(options.Username, options.Password, options.Target, tmpFileName)
		cmd = bmp.NewTargetCommand(options, fakeBmpClient)
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("NewTargetCommand", func() {
		It("create new TargetCommand", func() {
			Expect(cmd).ToNot(BeNil())

			cmd2 := bmp.NewTargetCommand(options, fakeBmpClient)
			Expect(cmd2).ToNot(BeNil())
			Expect(cmd2).To(Equal(cmd))
		})
	})

	Describe("#Name", func() {
		It("returns the name of a TargetCommand", func() {
			Expect(cmd.Name()).To(Equal("target"))
		})
	})

	Describe("#Description", func() {
		It("returns the description of a TargetCommand", func() {
			Expect(cmd.Description()).To(Equal("Set the URL of Bare Metal Provision Server"))
		})
	})

	Describe("#Usage", func() {
		It("returns the usage text of a TargetCommand", func() {
			Expect(cmd.Usage()).To(Equal("bmp target http://url.to.bmp.server"))
		})
	})

	Describe("#Options", func() {
		It("returns the options of a TargetCommand", func() {
			Expect(cmds.EqualOptions(cmd.Options(), options)).To(BeTrue())
		})
	})

	Describe("#Validate", func() {
		Context("when the options includes a URL", func() {
			It("returns true on validation", func() {
				validate, err := cmd.Validate()
				Expect(validate).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the options includes an empty target URL", func() {
			BeforeEach(func() {
				options.Target = ""
				cmd = bmp.NewTargetCommand(options, fakeBmpClient)
			})

			It("returns false on validation", func() {
				validate, err := cmd.Validate()
				Expect(validate).To(BeFalse())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the options includes bad target URL", func() {
			var badCmd1, badCmd2 cmds.Command

			BeforeEach(func() {
				options.Target = "bad-url"
				badCmd1 = bmp.NewTargetCommand(options, fakeBmpClient)

				options.Target = "..."
				badCmd2 = bmp.NewTargetCommand(options, fakeBmpClient)
			})

			It("returns false on validation", func() {
				validate, err := badCmd1.Validate()
				Expect(validate).To(BeFalse())
				Expect(err).To(HaveOccurred())

				validate, err = badCmd2.Validate()
				Expect(validate).To(BeFalse())
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("#Execute", func() {
		It("executes a good TargetCommand", func() {
			rc, err := cmd.Execute(args)

			fmt.Printf("====> err: %#v\n", err)

			Expect(rc).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())

			c := config.NewConfig(fakeBmpClient.ConfigPath())
			configInfo, err := c.LoadConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(configInfo.TargetUrl).To(Equal("http://fake.url"))
		})
	})
})
