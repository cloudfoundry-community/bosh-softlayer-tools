package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	main "github.com/cloudfoundry-community/bosh-softlayer-tools/main"
)

var _ = Describe("config", func() {
	var (
		config main.Config
		err    error
	)

	BeforeEach(func() {
		config = main.NewConfig("fake-path")
	})

	Describe("#GetPath", func() {
		Context("when a path is path to config", func() {
			It("returns the fake-path config path", func() {
				Expect(config.GetPath()).To(Equal("fake-path"))
			})
		})

		Context("when no path is set", func() {
			BeforeEach(func() {
				It("returns the default config path", func() {
					Expect(config.GetPath()).To(Equal(main.CONFIG_PATH))
				})
			})
		})
	})
})
