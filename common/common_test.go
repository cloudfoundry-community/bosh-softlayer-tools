package common_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	config "github.com/cloudfoundry-community/bosh-softlayer-tools/config"
)

var _ = Describe("LightStemcellCmd", func() {
	var (
		err                 error
		tmpDir, tmpFileName string
	)

	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir("", "bosh-softlayer-tools")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	XContext("CreateFile", func() {
		It("creates a file", func() {
			Fail("implement me!")
		})
	})

	XContext("CreateTarball", func() {
		It("creates a tarball", func() {
			Fail("implement me!")
		})
	})

	XContext("CreateBmpClient", func() {
		It("creates a BMP client", func() {
			Fail("implement me!")
		})
	})

	Context("#CreateConfig", func() {
		BeforeEach(func() {
			tmpFile, err := ioutil.TempFile(tmpDir, ".bmp_config")
			Expect(err).ToNot(HaveOccurred())

			tmpFileName = tmpFile.Name()

			c := config.NewConfig(tmpFileName)
			err = c.SaveConfig(config.ConfigInfo{})
			Expect(err).ToNot(HaveOccurred())

			configInfo, err := c.LoadConfig()
			Expect(err).ToNot(HaveOccurred())

			Expect(configInfo.Username).To(Equal(""))
			Expect(configInfo.Password).To(Equal(""))
			Expect(configInfo.TargetUrl).To(Equal(""))
		})

		It("creates config with data", func() {
			c := config.NewConfig(tmpFileName)
			Expect(err).ToNot(HaveOccurred())

			configInfo, err := c.LoadConfig()
			Expect(err).ToNot(HaveOccurred())

			configInfo.Username = "fake-username"
			configInfo.Password = "fake-password"
			configInfo.TargetUrl = "fake-target-url"

			err = c.SaveConfig(configInfo)
			Expect(err).ToNot(HaveOccurred())

			configInfo, err = c.LoadConfig()
			Expect(err).ToNot(HaveOccurred())

			Expect(configInfo.Username).To(Equal("fake-username"))
			Expect(configInfo.Password).To(Equal("fake-password"))
			Expect(configInfo.TargetUrl).To(Equal("fake-target-url"))
		})
	})
})
