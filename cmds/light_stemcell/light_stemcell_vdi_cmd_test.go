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

	testhelpers "github.com/maximilien/softlayer-go/test_helpers"
)

var _ = Describe("LightStemcellVDICmd", func() {
	var (
		err error

		fakeClient *slclientfakes.FakeSoftLayerClient

		options common.Options

		lightStemcellsPath, stemcellInfoFilePath string
		lightStemcellInfo                        LightStemcellInfo

		cmd *LightStemcellVDICmd
	)

	BeforeEach(func() {
		username := os.Getenv("SL_USERNAME")
		Expect(username).ToNot(Equal(""), "Missing SL_USERNAME environment variables")

		apiKey := os.Getenv("SL_API_KEY")
		Expect(apiKey).ToNot(Equal(""), "Missing SL_API_KEY environment variables")

		fakeClient = slclientfakes.NewFakeSoftLayerClient(username, apiKey)

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
			LightStemcellTypeFlag:    "VDI",
			InfrastructureFlag:       "fake-infrastructure",
			HypervisorFlag:           "fake-hypervisor",
			OsNameFlag:               "fake-os-name",
			LightStemcellPathFlag:    lightStemcellsPath,
			StemcellInfoFilenameFlag: stemcellInfoFilename,
		}

		cmd = NewLightStemcellVDICmd(options, fakeClient)
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
				cmd = NewLightStemcellVDICmd(options, fakeClient)
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
						Expect(err).To(HaveOccurred())
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
				fakeClient.FakeHttpClient.DoRawHttpRequestResponse, err = common.ReadJsonTestFixtures("../..", "softlayer", "SoftLayer_Virtual_Disk_Image_getObject.json")
				Expect(err).ToNot(HaveOccurred())
			})

			Describe("#Run", func() {
				It("runs the command", func() {
					err = cmd.Run()
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})

	Context("Light Stemcell Command interface", func() {
		BeforeEach(func() {
			fakeClient.FakeHttpClient.DoRawHttpRequestResponse, err = common.ReadJsonTestFixtures("../..", "softlayer", "SoftLayer_Virtual_Disk_Image_getObject.json")
			Expect(err).ToNot(HaveOccurred())
		})

		It("#NewLightStemcellCmd", func() {
			cmd := NewLightStemcellVDICmd(options, fakeClient)
			Expect(cmd.GetStemcellPath()).To(Equal(lightStemcellsPath))
		})

		It("#GenerateStemcellName", func() {
			name := GenerateStemcellName(lightStemcellInfo)
			Expect(name).ToNot(Equal(""))
			Expect(name).To(Equal("bosh-fake-infrastructure-fake-hypervisor-fake-os-name-go_agent"))
		})

		It("#Create", func() {
			lightStemcellPath, err := cmd.Create(1234567)
			Expect(err).ToNot(HaveOccurred())
			Expect(lightStemcellPath).ToNot(Equal(""), "the light stemcell path cannot be empty")
			Expect(testhelpers.FileExists(lightStemcellPath)).To(BeTrue())
		})

		It("#GetStemcellPath", func() {
			_, err := cmd.Create(1234567)
			Expect(err).ToNot(HaveOccurred())
			Expect(cmd.GetStemcellPath()).To(Equal(lightStemcellsPath))
		})
	})
})
