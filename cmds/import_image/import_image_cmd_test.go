package import_image_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/maximilien/bosh-softlayer-stemcells/cmds/import_image"

	cmds "github.com/maximilien/bosh-softlayer-stemcells/cmds"
	common "github.com/maximilien/bosh-softlayer-stemcells/common"
)

var _ = Describe("import-image command", func() {
	var (
		err     error
		cmd     cmds.CommandInterface
		options common.Options
	)

	Describe("#Options", func() {
		BeforeEach(func() {
			options = common.Options{
				NameFlag:      "fake-name",
				NoteFlag:      "fake-note",
				OsRefCodeFlag: "fake-os-ref-code",
				UriFlag:       "fake-uri",
			}
			cmd, err = NewImportImageCmd(options)
			Expect(err).ToNot(HaveOccurred())
		})

		It("contains a non-nil options", func() {
			Expect(cmd.Options()).ToNot(BeNil())

			Expect(cmd.Options().NameFlag).To(Equal("fake-name"))
			Expect(cmd.Options().NoteFlag).To(Equal("fake-note"))
			Expect(cmd.Options().OsRefCodeFlag).To(Equal("fake-os-ref-code"))
			Expect(cmd.Options().UriFlag).To(Equal("fake-uri"))
		})
	})
})
