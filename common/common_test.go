package common_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-community/bosh-softlayer-tools/common"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/config"
)

var _ = Describe("common", func() {
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

	XContext("CreateTarball", func() {
		It("creates a tarball", func() {
			Fail("implement me!")
		})
	})

	Context("ReadJsonTestFixtures", func() {
		type Test struct {
			Test string `json:"test"`
		}

		It("reads the test_fixtures/test/test.json", func() {
			contents, err := common.ReadJsonTestFixtures("..", "test", "test.json")
			Expect(err).NotTo(HaveOccurred())

			test := Test{}
			err = json.Unmarshal(contents, &test)
			Expect(err).NotTo(HaveOccurred())
			Expect(test.Test).To(Equal("test"))
		})
	})

	Context("CreateBmpClient", func() {
		var (
			currentUser    *user.User
			configFileName string
		)

		BeforeEach(func() {
			currentUser, err = user.Current()
			Expect(err).NotTo(HaveOccurred())

			configFileName = filepath.Join(currentUser.HomeDir, config.CONFIG_FILE_NAME)
		})

		Context("when current user has a .bmp_config", func() {
			BeforeEach(func() {
				_, err = ioutil.ReadFile(configFileName)
				if err != nil {
					configContents := []byte(`{
							"Username": "",
							"Password": "",
							"TargetUrl": ""
						}`)
					err = ioutil.WriteFile(configFileName, configContents, 0666)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			It("creates a BMP client", func() {
				bmpClient, err := common.CreateBmpClient()
				Expect(err).NotTo(HaveOccurred())
				Expect(bmpClient).ToNot(BeNil())
			})
		})

		Context("when current user does not have a .bmp_config", func() {
			var (
				tmpFileName string
				err         error
			)

			BeforeEach(func() {
				_, err = os.Stat(configFileName)
				if os.IsNotExist(err) == false {
					tmpFile, err := ioutil.TempFile("", ".bmp_config")
					Expect(err).NotTo(HaveOccurred())
					tmpFileName = tmpFile.Name()

					contents, err := ioutil.ReadFile(configFileName)
					Expect(err).NotTo(HaveOccurred())

					err = ioutil.WriteFile(tmpFileName, contents, 0666)
					Expect(err).NotTo(HaveOccurred())

					err = os.Remove(configFileName)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			AfterEach(func() {
				_, err := os.Stat(tmpFileName)
				if os.IsNotExist(err) == false {
					contents, err := ioutil.ReadFile(tmpFileName)
					Expect(err).NotTo(HaveOccurred())

					err = ioutil.WriteFile(configFileName, contents, 0666)
					Expect(err).NotTo(HaveOccurred())

					err = os.Remove(tmpFileName)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			It("fails to create a BMP client", func() {
				_, err = common.CreateBmpClient()
				Expect(err).To(HaveOccurred())
			})
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
