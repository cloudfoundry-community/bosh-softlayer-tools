package common_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LightStemcellCmd", func() {
	var (
		err    error
		tmpDir string
	)

	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir("", "bosh-softlayer-stemcells")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})
})
