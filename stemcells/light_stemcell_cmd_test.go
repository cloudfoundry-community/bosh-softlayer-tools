package stemcells_test

import (
	"os"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/maximilien/bosh-softlayer-stemcells/stemcells"

	common "github.com/maximilien/bosh-softlayer-stemcells/common"

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

		stemcellsPath string
		lightStemcellCmd LightStemcellCmd
	)

	BeforeEach(func() {
		fakeClient = slclientfakes.NewFakeSoftLayerClient("fake-username", "fake-api-key")

		accountService, err = testhelpers.CreateAccountService()
		Expect(err).ToNot(HaveOccurred())

		stemcellsPath, err = ioutil.TempDir("", "bosh-softlayer-stemcells")
		Expect(err).ToNot(HaveOccurred())

		lightStemcellCmd = NewLightStemcellCmd(stemcellsPath, fakeClient)
	})

	AfterEach(func() {
		err = os.RemoveAll(stemcellsPath)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Light Stemcell Command", func() {
		BeforeEach(func() {
			fakeClient.DoRawHttpRequestResponse, err = slgocommon.ReadJsonTestFixtures(".", "SoftLayer_Virtual_Disk_Image_getObject.json")
			Expect(err).ToNot(HaveOccurred())
		})

		It("#NewLightStemcellCmd", func() {
			cmd := NewLightStemcellCmd(stemcellsPath, fakeClient)
			Expect(cmd.GetStemcellsPath()).To(Equal(stemcellsPath))
		})

		XIt("#Create", func() {
			lightStemcellPath, err := lightStemcellCmd.Create(4868344)
			Expect(err).ToNot(HaveOccurred())
			Expect(lightStemcellPath).ToNot(Equal(""), "the light stemcell path cannot be empty")
			Expect(testhelpers.FileExists(lightStemcellPath)).To(BeTrue())
		})

		It("#GetStemcellsPath", func() {
			_, err := lightStemcellCmd.Create(4868344)
			Expect(err).ToNot(HaveOccurred())
			Expect(lightStemcellCmd.GetStemcellsPath()).To(Equal(stemcellsPath))
		})

		It("#GetLightStemcellMF", func() {
			expectedLightStemcellMF, err := common.LoadLightStemcellMF()
			Expect(err).ToNot(HaveOccurred())

			_, err := lightStemcellCmd.Create(4868344)
			Expect(err).ToNot(HaveOccurred())
			Expect(lightStemcellCmd.GetLightStemcellMF()).To(Equal(expectedLightStemcellMF))
		})

	})
})
