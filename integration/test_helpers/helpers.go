package test_helpers

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func RunBmp(args ...string) *Session {
	return RunCommand(BmpExec, args...)
}

func RunStemcells(args ...string) *Session {
	return RunCommand(StemcellsExec, args...)
}

func RunCommand(cmd string, args ...string) *Session {
	command := exec.Command(cmd, args...)

	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	session.Wait()

	return session
}
