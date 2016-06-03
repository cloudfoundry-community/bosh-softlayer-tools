package bmp_test

import (
	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("`$bmp update-state` integration tests", func() {
	var session *Session

	BeforeEach(func() {
		session = RunBmpTarget()
		Expect(session.ExitCode()).To(Equal(0))

		session = RunBmpLogin()
		Expect(session.ExitCode()).To(Equal(0))
	})

	Describe("$bmp update-state", func() {
		BeforeEach(func() {
			session = RunBmp("update-state", "--server", "fake-id", "--state", "bm.state.new")
		})

		It("returns 1", func() {
			Expect(session.ExitCode()).To(Equal(1))
		})

		It("prints the error message", func() {
			Expect(session.Wait().Out.Contents()).Should(ContainSubstring("Update Cancelled!"))
		})
	})
})
