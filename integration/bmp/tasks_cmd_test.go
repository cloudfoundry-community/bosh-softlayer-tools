package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
)

var _ = Describe("`$bmp tasks` integration tests", func() {
	var session *Session

	BeforeEach(func() {
		session = RunBmpTarget()
		Expect(session.ExitCode()).To(Equal(0))

		session = RunBmpLogin()
		Expect(session.ExitCode()).To(Equal(0))
	})

	Describe("$bmp tasks", func() {
		Context("execute bmp tasks without --latest", func() {
			BeforeEach(func() {
				session = RunBmp("tasks")
			})

			It("returns 0", func() {
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("prints the BMP tasks information", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("Tasks total"))
			})
		})

		Context("execute bmp tasks with --latest", func() {
			BeforeEach(func() {
				session = RunBmp("tasks", "--latest=10")
			})

			It("returns 0", func() {
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("prints the BMP tasks information", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("Tasks total"))
			})
		})

	})
})
