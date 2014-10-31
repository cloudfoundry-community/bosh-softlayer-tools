package stemcells_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/maximilien/bosh-softlayer-stemcells/stemcells"

	slclientfakes "github.com/maximilien/softlayer-go/client/fakes"

	slgocommon "github.com/maximilien/softlayer-go/common"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
	testhelpers "github.com/maximilien/softlayer-go/test_helpers"
)

var _ = Describe("LightStemcellCmd", func() {
	var (
		err error

		fakeClient *slclientfakes.FakeSoftLayerClient

		accountService softlayer.SoftLayer_Account_Service

		lightStemcellsPath string
		lightStemcellInfo  LightStemcellInfo
		lightStemcellCmd   LightStemcellCmd
	)

	BeforeEach(func() {
		username := os.Getenv("SL_USERNAME")
		Expect(username).ToNot(Equal(""))

		apiKey := os.Getenv("SL_API_KEY")
		Expect(apiKey).ToNot(Equal(""))

		fakeClient = slclientfakes.NewFakeSoftLayerClient(username, apiKey)

		accountService, err = testhelpers.CreateAccountService()
		Expect(err).ToNot(HaveOccurred())

		lightStemcellsPath, err = ioutil.TempDir("", "bosh-softlayer-stemcells")
		Expect(err).ToNot(HaveOccurred())

		lightStemcellInfo = LightStemcellInfo{
			Infrastructure: "fake-infrastructure",
			Architecture:   "fake-architecture",
			RootDeviceName: "fake-root-device-name",

			Version:    "fake-version",
			Hypervisor: "fake-hypervisor",
			OsName:     "fake-os-name",
		}

		lightStemcellCmd = NewLightStemcellCmd(lightStemcellsPath, lightStemcellInfo, fakeClient)
	})

	AfterEach(func() {
		err = os.RemoveAll(lightStemcellsPath)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Light Stemcell Command", func() {
		BeforeEach(func() {
			fakeClient.DoRawHttpRequestResponse, err = slgocommon.ReadJsonTestFixtures(".", "SoftLayer_Virtual_Disk_Image_getObject.json")
			Expect(err).ToNot(HaveOccurred())
		})

		It("#NewLightStemcellCmd", func() {
			cmd := NewLightStemcellCmd(lightStemcellsPath, lightStemcellInfo, fakeClient)
			Expect(cmd.GetStemcellsPath()).To(Equal(lightStemcellsPath))
		})

		It("#GenerateStemcellName", func() {
			name := lightStemcellCmd.GenerateStemcellName(lightStemcellInfo)
			Expect(name).ToNot(Equal(""))
			Expect(name).To(Equal("light-bosh-stemcell-fake-version-fake-infrastructure-fake-hypervisor-fake-os-name-go_agent"))
		})

		It("#Create", func() {
			lightStemcellPath, err := lightStemcellCmd.Create(1234567)
			Expect(err).ToNot(HaveOccurred())
			Expect(lightStemcellPath).ToNot(Equal(""), "the light stemcell path cannot be empty")
			Expect(testhelpers.FileExists(lightStemcellPath)).To(BeTrue())
		})

		It("#GetStemcellsPath", func() {
			_, err := lightStemcellCmd.Create(1234567)
			Expect(err).ToNot(HaveOccurred())
			Expect(lightStemcellCmd.GetStemcellsPath()).To(Equal(lightStemcellsPath))
		})
	})
})
