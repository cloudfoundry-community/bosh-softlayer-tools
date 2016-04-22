package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
)

var _ = Describe("`$bmp status` integration tests", func() {
	var session *Session

	Describe("$bmp status", func() {
		BeforeEach(func() {
			session = RunBmp("status")
		})

		It("returns 0", func() {
			Expect(session.ExitCode()).To(Equal(0))
		})

		It("prints a the BMP server status", func() {
			Expect(session).To(Say("BMP status"))
		})
	})
})
