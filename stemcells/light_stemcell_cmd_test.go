package stemcells_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	slclientfakes "github.com/maximilien/softlayer-go/client/fakes"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
	testhelpers "github.com/maximilien/softlayer-go/test_helpers"
)

var _ = Describe("LightStemcellCmd", func() {
	var (
		err error

		client *slclientfakes.FakeSoftLayerClient

		accountService softlayer.SoftLayer_Account_Service
	)

	BeforeEach(func() {
		client = slclientfakes.NewFakeSoftLayerClient("fake-username", "fake-api-key")

		accountService, err = testhelpers.CreateAccountService()
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Light Stemcell Command", func() {
		It("#NewLightStemcellCmd", func() {
			Fail("implement me!")
		})

		It("#Create", func() {
			Fail("implement me!")
		})
	})
})
