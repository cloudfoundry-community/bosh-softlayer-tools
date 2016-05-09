package bmp_test

import (
	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("`$bmp sl` integration tests", func() {
	var session *Session

	BeforeEach(func() {
		session = RunBmpTarget()
		Expect(session.ExitCode()).To(Equal(0))

		session = RunBmpLogin()
		Expect(session.ExitCode()).To(Equal(0))
	})

	Describe("$bmp sl", func() {
		Context("execute --packages", func() {
			BeforeEach(func() {
				session = RunBmp("sl", "--packages")
			})

			It("returns 0", func() {
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("prints the package information", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("Packages total"))
			})
		})
		//TODO: Add test cases for --package-options
	})
})
