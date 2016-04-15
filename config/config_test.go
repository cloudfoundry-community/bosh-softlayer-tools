package config_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	config "github.com/cloudfoundry-community/bosh-softlayer-tools/config"
)

var _ = Describe("config", func() {
	var (
		c                       config.Config
		err                     error
		tmpDirName, tmpFileName string
	)

	BeforeEach(func() {
		c = config.NewConfig("fake-path")
	})

	Describe("#Path", func() {
		Context("when a path is path to config", func() {
			It("returns the fake-path config path", func() {
				Expect(c.Path()).To(Equal("fake-path"))
			})
		})

		Context("when no path is set", func() {
			BeforeEach(func() {
				c = config.NewConfig("")
			})

			It("returns the default config path", func() {
				Expect(c.Path()).To(Equal(config.CONFIG_PATH))
			})
		})
	})

	Describe("#LoadConfig", func() {
		AfterEach(func() {
			_, err := os.Stat(tmpDirName)
			if os.IsExist(err) {
				err = os.RemoveAll(tmpDirName)
				Expect(err).To(HaveOccurred())
			}
		})

		Context("when config file does not exist", func() {
			It("returns error", func() {
				_, err = c.LoadConfig()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when config file is not valid config file", func() {
			BeforeEach(func() {
				tmpDirName, err = ioutil.TempDir("", "config")
				Expect(err).ToNot(HaveOccurred())

				tmpFileName = filepath.Join(tmpDirName, "config")
				err = ioutil.WriteFile(tmpFileName, []byte{}, 0666)
				Expect(err).ToNot(HaveOccurred())

				c = config.NewConfig(tmpFileName)
			})

			It("returns error", func() {
				_, err = c.LoadConfig()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when config is valid", func() {
			BeforeEach(func() {
				tmpDirName, err = ioutil.TempDir("", "config")
				Expect(err).ToNot(HaveOccurred())

				tmpFileName = filepath.Join(tmpDirName, "config")
				err = ioutil.WriteFile(tmpFileName, []byte(
					`{
    "username": "fake-username",
	"password": "fake-password",
	"target_url": "http://fake.target.url"
}`), 0666)
				Expect(err).ToNot(HaveOccurred())

				c = config.NewConfig(tmpFileName)
			})

			It("returns the config object", func() {
				configInfo, err := c.LoadConfig()
				Expect(err).ToNot(HaveOccurred())

				Expect(configInfo.Username).To(Equal("fake-username"))
				Expect(configInfo.Password).To(Equal("fake-password"))
				Expect(configInfo.TargetUrl).To(Equal("http://fake.target.url"))
			})
		})
	})

	Describe("#SaveConfig", func() {
		BeforeEach(func() {
			tmpDirName, err = ioutil.TempDir("", "config")
			Expect(err).ToNot(HaveOccurred())

			tmpFileName = filepath.Join(tmpDirName, "config")
			err = ioutil.WriteFile(tmpFileName, []byte(
				`{
    "username": "",
	"password": "",
	"target_url": ""
}`), 0666)
			Expect(err).ToNot(HaveOccurred())

			c = config.NewConfig(tmpFileName)
		})

		It("saves new content of config in config path", func() {
			configInfo := config.ConfigInfo{
				Username:  "fake-username",
				Password:  "fake-password",
				TargetUrl: "http://fake.target.url",
			}

			err = c.SaveConfig(configInfo)
			Expect(err).ToNot(HaveOccurred())

			configFileContents, err := ioutil.ReadFile(tmpFileName)
			Expect(err).ToNot(HaveOccurred())

			readConfigInfo := config.ConfigInfo{}
			err = json.Unmarshal(configFileContents, &readConfigInfo)
			Expect(err).ToNot(HaveOccurred())

			Expect(readConfigInfo.Username).To(Equal("fake-username"))
			Expect(readConfigInfo.Password).To(Equal("fake-password"))
			Expect(readConfigInfo.TargetUrl).To(Equal("http://fake.target.url"))
		})
	})
})
