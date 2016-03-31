package clients_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"

	slclientfakes "github.com/maximilien/softlayer-go/client/fakes"
)

var _ = Describe("BMP client", func() {

	var (
		err            error
		bmpClient      clients.BmpClient
		fakeHttpClient *slclientfakes.FakeHttpClient
	)

	BeforeEach(func() {
		fakeHttpClient = slclientfakes.NewFakeHttpClient("fake-username", "fake-password")
		Expect(fakeHttpClient).ToNot(BeNil())

		fakeHttpClient.DoRawHttpRequestResponse, err = common.ReadJsonTestFixtures("..", "bmp", "Info.json")
		Expect(err).ToNot(HaveOccurred())

		bmpClient = clients.NewBmpClient("fake-username", "fake-password", "http://fake.url.com", fakeHttpClient)
		Expect(bmpClient).ToNot(BeNil())
	})

	Describe("#Info", func() {
		It("returns BMP server info", func() {
			info, err := bmpClient.Info()
			Expect(err).ToNot(HaveOccurred())

			Expect(info).To(Equal(clients.InfoResponse{
				Status: 0,
				Data: clients.DataInfo{
					Name:    "fake-name",
					Version: "fake-version"}}))
		})

		It("fails when BMP server fails", func() {
			fakeHttpClient.DoRawHttpRequestError = errors.New("fake-error")

			_, err := bmpClient.Info()
			Expect(err).To(HaveOccurred())
		})
	})
})
