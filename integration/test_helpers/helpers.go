package test_helpers

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func RunBmp(args ...string) *Session {
	return RunCommand(BmpExec, args...)
}

func RunBmpTarget() *Session {
	TargetURL, err := GetTargetURL()
	Expect(err).ToNot(HaveOccurred())

	session := RunBmp("target", "--target", TargetURL)
	return session
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

func GetTargetURL() (string, error) {
	TargetURL := os.Getenv("TARGET_URL")
	if TargetURL == "" {
		return "", errors.New("TARGET_URL environment must be set")
	}

	_, err := url.ParseRequestURI(TargetURL)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s is not a valid URL", TargetURL))
	}

	return TargetURL, nil
}
