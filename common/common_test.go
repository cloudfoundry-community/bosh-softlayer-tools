package common_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	common "github.com/maximilien/bosh-softlayer-stemcells/common"

	stemcells "github.com/maximilien/bosh-softlayer-stemcells/stemcells"
)

var _ = Describe("LightStemcellCmd", func() {
	var (
		err               error
		lightStemcellMF   stemcells.LightStemcellMF
		lightStemcellPath string
		tmpDir            string
	)

	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir("", "bosh-softlayer-stemcells")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	Context("common.LoadStemcellMF", func() {
		BeforeEach(func() {
			lightStemcellPath = filepath.Join(tmpDir, "fake-stemcell.tgz")
		})

		It("loads and parses a new StemcellMF object from json file", func() {
			lsMF, err := common.LoadLightStemcellMF(lightStemcellPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(lsMF).To(Equal(lightStemcellMF))
		})
	})
})
