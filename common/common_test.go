package common_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = XDescribe("LightStemcellCmd", func() {
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

	Context("CreateFile", func() {
		It("creates a file", func() {
			Fail("implement me!")
		})
	})

	Context("CreateTarball", func() {
		It("creates a tarball", func() {
			Fail("implement me!")
		})
	})

	Context("CreateBmpClient", func() {
		It("creates a BMP client", func() {
			Fail("implement me!")
		})
	})
})
