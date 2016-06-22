package light_stemcell_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/light_stemcell"
)

var _ = Describe("light_stemcell_helper", func() {
	var (
		err               error
		tmpDir            string
		lightStemcellInfo LightStemcellInfo
		lightStemcellMF   LightStemcellMF
	)

	BeforeEach(func() {
		tmpDir, err = ioutil.TempDir("", "bosh-softlayer-stemcells")
		Expect(err).ToNot(HaveOccurred())

		lightStemcellInfo = LightStemcellInfo{
			Infrastructure: "fake-infrastructure",
			Architecture:   "fake-architecture",
			RootDeviceName: "fake-root-device-name",

			Version:    "fake-version",
			Hypervisor: "fake-hypervisor",
			OsName:     "fake-os-name",
		}

		lightStemcellMF = LightStemcellMF{
			Name:         "fake-name",
			Version:      "fake-version",
			BoshProtocol: 1,
			Sha1:         "fake-sha1",
			CloudProperties: CloudProperties{
				Infrastructure: "fake-infrastructure",
				Architecture:   "fake-architecture",
				RootDeviceName: "/fake/root/device",

				VirtualDiskImageId:   12345,
				VirtualDiskImageUuid: "fake-uuid",
				DatacenterName:       "fake-datacenter-name",
			},
		}
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	It("GenerateStemcellName", func() {
		name := GenerateStemcellName(lightStemcellInfo)
		Expect(name).To(Equal("bosh-fake-infrastructure-fake-hypervisor-fake-os-name-go_agent"))
	})

	It("GenerateStemcellFileName", func() {
		name := GenerateStemcellFilelName(lightStemcellInfo)
		Expect(name).To(Equal("light-bosh-stemcell-fake-version-fake-infrastructure-fake-hypervisor-fake-os-name-go_agent"))
	})

	It("GenerateLightStemcellTarball", func() {
		tarballPath, err := GenerateLightStemcellTarball(lightStemcellMF, lightStemcellInfo, tmpDir)
		Expect(err).ToNot(HaveOccurred())
		Expect(tarballPath).ToNot(Equal(""))
		Expect(tarballPath).To(ContainSubstring(tmpDir))

		stat, err := os.Stat(tarballPath)
		Expect(err).ToNot(HaveOccurred())
		Expect(stat.Name()).To(Equal(fmt.Sprintf("%s.tgz", GenerateStemcellFilelName(lightStemcellInfo))))
		Expect(stat.Size()).To(BeNumerically(">", 0))
		Expect(stat.IsDir()).To(BeFalse())
	})

	It("GenerateManifestMFBytesJSON", func() {
		bytes, err := GenerateManifestMFBytesJSON(lightStemcellMF)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(bytes)).To(BeNumerically(">", 0))
	})

	It("GenerateManifestMFBytesYAML", func() {
		bytes, err := GenerateManifestMFBytesYAML(lightStemcellMF)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(bytes)).To(BeNumerically(">", 0))
	})
})
