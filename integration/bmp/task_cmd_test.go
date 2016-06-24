package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
)

var _ = Describe("`$bmp task --task_id` integration tests", func() {
	var session *Session

	BeforeEach(func() {
		session = RunBmpTarget()
		Expect(session.ExitCode()).To(Equal(0))

		session = RunBmpLogin()
		Expect(session.ExitCode()).To(Equal(0))
	})

	Describe("$bmp task --task_id", func() {
		Context("when execute bmp task with default event level", func() {
			BeforeEach(func() {
				session = RunBmp("task", "--task_id=1")
			})

			It("returns 0", func() {
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("prints task output", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("Task output for ID 1 with event level"))
			})
		})

		Context("when execute bmp task with debug level", func() {
			BeforeEach(func() {
				session = RunBmp("task", "--task_id=1", "--debug")
			})

			It("returns 0", func() {
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("prints task output", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("Task output for ID 1 with debug level"))
			})
		})

		Context("when execute bmp task with json format and event level", func() {
			BeforeEach(func() {
				session = RunBmp("task", "--task_id=1", "--json")
			})

			It("returns 0", func() {
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("prints task output", func() {
				Expect(session.Wait().Out.Contents()).To(ContainSubstring("Task output for ID 1 with event level"))
			})
		})
	})
})
