package bmp_test

import (
	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("`$bmp create-baremetals` integration tests", func() {
	var session *Session

	BeforeEach(func() {
		session = RunBmpTarget()
		Expect(session.ExitCode()).To(Equal(0))

		session = RunBmpLogin()
		Expect(session.ExitCode()).To(Equal(0))
	})

	Describe("$bmp create-baremetals", func() {
		Context("when executes create-baremetals", func() {
			BeforeEach(func() {
				deployment, err := GetDeployment()
				Expect(err).ToNot(HaveOccurred())

				session = RunBmp("create-baremetals", "--deployment", deployment)
			})

			It("returns 0", func() {
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("returns the task id", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("Run command: bmp task"))
			})
		})

		Context("when executes create-baremetals --dryrun", func() {
			BeforeEach(func() {
				deployment, err := GetDeployment()
				Expect(err).ToNot(HaveOccurred())

				session = RunBmp("create-baremetals", "--deployment", deployment, "--dryrun")
			})

			It("returns 0", func() {
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("returns the task id", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("Run command: bmp task"))
			})
		})
	})
})
