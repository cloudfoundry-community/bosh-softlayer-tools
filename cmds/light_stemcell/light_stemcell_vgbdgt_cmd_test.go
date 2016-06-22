package light_stemcell_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/light_stemcell"

	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"

	slclientfakes "github.com/maximilien/softlayer-go/client/fakes"

	softlayer "github.com/maximilien/softlayer-go/softlayer"
	testhelpers "github.com/maximilien/softlayer-go/test_helpers"
)

var _ = Describe("LightStemcellVGBDGTCmd", func() {
	var (
		err error

		fakeClient *slclientfakes.FakeSoftLayerClient

		options common.Options

		accountService softlayer.SoftLayer_Account_Service

		lightStemcellsPath, stemcellInfoFilePath string
		lightStemcellInfo                        LightStemcellInfo

		cmd *LightStemcellVGBDTGCmd
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

		stemcellInfoFilePath, err = ioutil.TempDir("", "bosh-softlayer-stemcells")
		Expect(err).ToNot(HaveOccurred())

		stemcellInfoFileContent := `{
			"id": 1234567,
			"uuid": "fake-uuid"
		}`

		stemcellInfoFilename := filepath.Join(stemcellInfoFilePath, "stemcell-info.json")
		err = ioutil.WriteFile(stemcellInfoFilename, []byte(stemcellInfoFileContent), 0644)
		Expect(err).ToNot(HaveOccurred())

		options = common.Options{
			NameFlag:                 "fake-name",
			VersionFlag:              "fake-version",
			LightStemcellTypeFlag:    "VGBDGT",
			InfrastructureFlag:       "fake-infrastructure",
			HypervisorFlag:           "fake-hypervisor",
			OsNameFlag:               "fake-os-name",
			LightStemcellPathFlag:    lightStemcellsPath,
			StemcellInfoFilenameFlag: stemcellInfoFilename,
		}

		cmd = NewLightStemcellVGBDGTCmd(options, fakeClient)
		lightStemcellInfo = cmd.GetLightStemcellInfo()
	})

	AfterEach(func() {
		err = os.RemoveAll(lightStemcellsPath)
		Expect(err).ToNot(HaveOccurred())

		err = os.RemoveAll(stemcellInfoFilePath)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Command interface methods", func() {
		Describe("#Options", func() {
			It("return the options struct", func() {
				options := cmd.Options()
				Expect(options).ToNot(BeNil())
				Expect(options.NameFlag).To(Equal("fake-name"))
				Expect(options.VersionFlag).To(Equal("fake-version"))

				//TODO more tests
			})
		})

		Describe("#CheckOptions", func() {
			JustBeforeEach(func() {
				cmd = NewLightStemcellVGBDGTCmd(options, fakeClient)
				Expect(cmd).ToNot(BeNil())
			})

			Context("when all required options are passed", func() {
				It("succeeds with no errors", func() {
					err = cmd.CheckOptions()
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("when one required option is missing", func() {
				Context("when Version is missing", func() {
					BeforeEach(func() {
						options.VersionFlag = ""
					})

					It("fails with warning that version is missing", func() {
						err = cmd.CheckOptions()
						Expect(err.Error()).To(ContainSubstring("must pass a version"))
					})
				})

				Context("when Stemcell Info Filename is missing", func() {
					BeforeEach(func() {
						options.StemcellInfoFilenameFlag = ""
					})

					It("fails with warning that stemcell-info.json filename flag is missing", func() {
						err = cmd.CheckOptions()
						Expect(err.Error()).To(ContainSubstring("must pass a path to stemcell-info.json"))
					})
				})
			})

			Context("when all required options are missing", func() {
				BeforeEach(func() {
					options.VersionFlag = ""
				})

				It("fails with warning that version is missing", func() {
					err = cmd.CheckOptions()
					Expect(err.Error()).To(ContainSubstring("must pass a version"))
				})
			})
		})

		Context("running command", func() {
			BeforeEach(func() {
				readJsonTestFixtures(fakeClient)
			})

			Describe("#Run", func() {
				It("runs the command", func() {
					err = cmd.Run()
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})

	Context("Light Stemcell Command interface methods", func() {
		BeforeEach(func() {
			readJsonTestFixtures(fakeClient)
		})

		Describe("#NewLightStemcellCmd", func() {
			It("creates a new light stemcell command", func() {
				cmd = NewLightStemcellVGBDGTCmd(options, fakeClient)
				Expect(cmd.GetStemcellPath()).To(Equal(lightStemcellsPath))
			})
		})

		Describe("#GenerateStemcellName", func() {
			It("generates a stemcell name", func() {
				name := GenerateStemcellName(lightStemcellInfo)
				Expect(name).ToNot(Equal(""))
				Expect(name).To(Equal("bosh-fake-infrastructure-fake-hypervisor-fake-os-name-go_agent"))
			})
		})

		Describe("#Create", func() {
			It("creates the light stemcell file", func() {
				lightStemcellPath, err := cmd.Create(1234567)
				Expect(err).ToNot(HaveOccurred())
				Expect(lightStemcellPath).ToNot(Equal(""), "the light stemcell path cannot be empty")
				Expect(testhelpers.FileExists(lightStemcellPath)).To(BeTrue())
			})
		})

		Describe("#GetStemcellPath", func() {
			It("returns the light stemcell path", func() {
				_, err := cmd.Create(1234567)
				Expect(err).ToNot(HaveOccurred())
				Expect(cmd.GetStemcellPath()).To(Equal(lightStemcellsPath))
			})
		})
	})
})

func readJsonTestFixtures(fakeClient *slclientfakes.FakeSoftLayerClient) {
	blockDeviceTemplateGroups, err := common.ReadJsonTestFixtures("../..", "softlayer", "SoftLayer_Account_Service_getBlockDeviceTemplateGroups.json")
	Expect(err).ToNot(HaveOccurred())

	getObject, err := common.ReadJsonTestFixtures("../..", "softlayer", "SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service_getObject.json")
	Expect(err).ToNot(HaveOccurred())

	getDatacenters, err := common.ReadJsonTestFixtures("../..", "softlayer", "SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service_getDatacenters.json")
	Expect(err).ToNot(HaveOccurred())

	fakeClient.FakeHttpClient.DoRawHttpRequestResponses = [][]byte{blockDeviceTemplateGroups, getObject, getDatacenters}
}
