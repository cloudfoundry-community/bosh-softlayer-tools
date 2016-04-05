package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	config "github.com/cloudfoundry-community/bosh-softlayer-tools/config"
)

var _ = Describe("config", func() {
	var (
		c config.Config
	)

	BeforeEach(func() {
		c = config.NewConfig("fake-path")
	})

	Describe("#GetPath", func() {
		Context("when a path is path to config", func() {
			It("returns the fake-path config path", func() {
				Expect(c.GetPath()).To(Equal("fake-path"))
			})
		})

		Context("when no path is set", func() {
			BeforeEach(func() {
				It("returns the default config path", func() {
					Expect(c.GetPath()).To(Equal(config.CONFIG_PATH))
				})
			})
		})
	})
})
