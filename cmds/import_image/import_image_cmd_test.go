package import_image_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/import_image"

	testhelpers "github.com/cloudfoundry-community/bosh-softlayer-tools/test_helpers"
	slclientfakes "github.com/maximilien/softlayer-go/client/fakes"
	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

var _ = Describe("import-image command", func() {
	var (
		err error

		fakeClient    *slclientfakes.FakeSoftLayerClient
		vgbdtgService softlayer.SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service

		cmd     cmds.CommandInterface
		options common.Options

		importImageCmd *ImportImageCmd
	)

	BeforeEach(func() {
		username := os.Getenv("SL_USERNAME")
		Expect(username).ToNot(Equal(""))

		apiKey := os.Getenv("SL_API_KEY")
		Expect(apiKey).ToNot(Equal(""))

		fakeClient = slclientfakes.NewFakeSoftLayerClient(username, apiKey)
		Expect(fakeClient).ToNot(BeNil())

		vgbdtgService, err = fakeClient.GetSoftLayer_Virtual_Guest_Block_Device_Template_Group_Service()
		Expect(err).ToNot(HaveOccurred())
		Expect(vgbdtgService).ToNot(BeNil())

		options = common.Options{
			NameFlag:       "fake-name",
			NoteFlag:       "fake-note",
			PublicFlag:     false,
			PublicNameFlag: "fake-public-name",
			PublicNoteFlag: "fake-public-note",
			OsRefCodeFlag:  "fake-os-ref-code",
			UriFlag:        "fake-uri",
		}

		importImageCmd, err = NewImportImageCmd(options, fakeClient)
		Expect(err).ToNot(HaveOccurred())
		Expect(importImageCmd).ToNot(BeNil())

		cmd = importImageCmd
	})

	Describe("#Options", func() {
		It("contains a non-nil options", func() {
			Expect(cmd.Options()).ToNot(BeNil())

			Expect(cmd.Options().NameFlag).To(Equal("fake-name"))
			Expect(cmd.Options().NoteFlag).To(Equal("fake-note"))
			Expect(cmd.Options().OsRefCodeFlag).To(Equal("fake-os-ref-code"))
			Expect(cmd.Options().UriFlag).To(Equal("fake-uri"))
		})
	})

	Describe("#CheckOptions", func() {
		JustBeforeEach(func() {
			importImageCmd, err = NewImportImageCmd(options, fakeClient)
			Expect(err).ToNot(HaveOccurred())
			Expect(importImageCmd).ToNot(BeNil())

			cmd = importImageCmd
		})

		Context("when all required options are passed", func() {
			It("succeeds with no errors", func() {
				err = cmd.CheckOptions()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when no required options are passed", func() {
			BeforeEach(func() {
				options.OsRefCodeFlag = ""
				options.UriFlag = ""
			})

			It("fails with error that operating system reference code is missing", func() {
				err = cmd.CheckOptions()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("must pass an OS ref code"))
			})
		})

		Context("when one required option is missing", func() {
			Context("when OsRefCode is missing", func() {
				BeforeEach(func() {
					options.OsRefCodeFlag = ""
				})

				It("fails with error that operating system reference code is missing", func() {
					err = cmd.CheckOptions()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("must pass an OS ref code"))
				})
			})

			Context("when Uri is missing", func() {
				BeforeEach(func() {
					options.UriFlag = ""
				})

				It("fails with error that the URI is missing", func() {
					err = cmd.CheckOptions()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("must pass a URI"))
				})
			})
		})
	})

	Describe("#Run", func() {
		var configuration sldatatypes.SoftLayer_Container_Virtual_Guest_Block_Device_Template_Configuration

		BeforeEach(func() {
			configuration = sldatatypes.SoftLayer_Container_Virtual_Guest_Block_Device_Template_Configuration{
				Name: "fake-configuration-name",
				Note: "fake-configuration-note",
				OperatingSystemReferenceCode: "fake-operating-system-reference-code",
				Uri: "swift://FakeObjectStorageAccountName>@fake-clusterName/fake-containerName/fake-fileName.vhd",
			}
		})

		It("creates a private VGDTG with UUID and ID", func() {
			fileNames := []string{
				"SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service_createFromExternalSource.json",
				"SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service_createFromExternalSource.json",
			}
			testhelpers.SetTestFixturesForFakeSoftLayerClient(fakeClient, fileNames)

			vgbdtGroup, err := vgbdtgService.CreateFromExternalSource(configuration)
			Expect(err).ToNot(HaveOccurred())

			Expect(vgbdtGroup.Id).To(Equal(211582))
			Expect(vgbdtGroup.GlobalIdentifier).To(Equal("fake-global-identifier"))

			err = cmd.Run()
			Expect(err).ToNot(HaveOccurred())

			Expect(importImageCmd.Uuid).ToNot(Equal(""))
			Expect(importImageCmd.Id).ToNot(Equal(""))
		})

		It("creates a public VGDTG with UUID and ID", func() {
			fileNames := []string{
				"SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service_createFromExternalSource.json",
				"SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service_createPublicArchiveTransaction.json",
			}
			testhelpers.SetTestFixturesForFakeSoftLayerClient(fakeClient, fileNames)

			options = common.Options{
				NameFlag:       "fake-name",
				NoteFlag:       "fake-note",
				PublicFlag:     true,
				PublicNameFlag: "fake-public-name",
				PublicNoteFlag: "fake-public-note",
				OsRefCodeFlag:  "fake-os-ref-code",
				UriFlag:        "fake-uri",
			}

			importImageCmd, err = NewImportImageCmd(options, fakeClient)
			Expect(err).ToNot(HaveOccurred())
			Expect(importImageCmd).ToNot(BeNil())

			cmd = importImageCmd

			err = cmd.Run()
			Expect(err).ToNot(HaveOccurred())

			Expect(importImageCmd.Uuid).To(Equal(""))
			Expect(importImageCmd.Id).ToNot(Equal(""))
		})
	})
})
