package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
)

var _ = Describe("`$bmp stemcells` integration tests", func() {
	var session *Session

	BeforeEach(func() {
		session = RunBmpTarget()
		Expect(session.ExitCode()).To(Equal(0))
	})

	Describe("$bmp stemcells", func() {
		BeforeEach(func() {
			session = RunBmp("stemcells")
		})

		It("returns 0", func() {
			Expect(session.ExitCode()).To(Equal(0))
		})

		It("prints the BMP stemcells information", func() {
			Expect(session.Wait().Out.Contents()).To(ContainSubstring("Stemcells total"))
		})
	})
})
