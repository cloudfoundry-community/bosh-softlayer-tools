package test_helpers

import (
	"errors"
	"os"
	"os/exec"

	config "github.com/cloudfoundry-community/bosh-softlayer-tools/config"
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

func RunBmpLogin() *Session {
	Username, Password, err := GetCredential()
	Expect(err).ToNot(HaveOccurred())

	session := RunBmp("login", "--username", Username, "--password", Password)
	Expect(session.ExitCode()).To(Equal(0))

	c := config.NewConfig("")
	configInfo, err := c.LoadConfig()
	Expect(configInfo.Username).To(Equal(Username))
	Expect(configInfo.Password).To(Equal(Password))

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

	return TargetURL, nil
}

func GetCredential() (string, string, error) {
	Username := os.Getenv("BMP_USERNAME")
	if Username == "" {
		return "", "", errors.New("BMP_USERNAME environment must be set")
	}

	Password := os.Getenv("BMP_PASSWORD")
	if Password == "" {
		return "", "", errors.New("BMP_PASSWORD environment must be set")
	}
	return Username, Password, nil
}
