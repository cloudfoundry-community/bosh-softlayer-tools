package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
)

var _ = Describe("`$bmp bms --deployment` integration tests", func() {
	var session *Session

	BeforeEach(func() {
		session = RunBmpTarget()
		Expect(session.ExitCode()).To(Equal(0))

		session = RunBmpLogin()
		Expect(session.ExitCode()).To(Equal(0))
	})

	Describe("$bmp bms --deployment", func() {
		BeforeEach(func() {
			deployment, _ := GetDeployment()
			session = RunBmp("bms", "--deployment", deployment)
		})

		It("returns 0", func() {
			Expect(session.ExitCode()).To(Equal(0))
		})

		It("prints baremetal server information", func() {
			Expect(session.Wait().Out.Contents()).Should(ContainSubstring("Baremetals total"))
		})
	})
})
