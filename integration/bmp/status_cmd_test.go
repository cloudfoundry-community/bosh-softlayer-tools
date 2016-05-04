package bmp_test

import (
	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("`$bmp status` integration tests", func() {
	var session *Session

	BeforeEach(func() {
		session = RunBmpTarget()
		Expect(session.ExitCode()).To(Equal(0))

		session = RunBmpLogin()
		Expect(session.ExitCode()).To(Equal(0))
	})

	Describe("$bmp status", func() {
		BeforeEach(func() {
			session = RunBmp("status")
		})

		It("returns 0", func() {
			Expect(session.ExitCode()).To(Equal(0))
		})

		It("prints a the BMP server status", func() {
			Expect(session.Wait().Out.Contents()).Should(ContainSubstring("Bluemix Provision Server"))
		})
	})
})
