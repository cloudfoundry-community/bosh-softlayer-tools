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

		bmpClient = clients.NewBmpClient("fake-username", "fake-password", "http://fake.url.com", fakeHttpClient)
		Expect(bmpClient).ToNot(BeNil())
	})

	Describe("#Info", func() {
		BeforeEach(func() {
			fakeHttpClient.DoRawHttpRequestResponse, err = common.ReadJsonTestFixtures("..", "bmp", "Info.json")
			Expect(err).ToNot(HaveOccurred())
		})

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

	Describe("#SlPackages", func() {
		BeforeEach(func() {
			fakeHttpClient.DoRawHttpRequestResponse, err = common.ReadJsonTestFixtures("..", "bmp", "SlPackages.json")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an array of DataPackage", func() {
			slPackageResponse, err := bmpClient.SlPackages()
			Expect(err).ToNot(HaveOccurred())

			Expect(len(slPackageResponse.Data)).To(Equal(2))
			Expect(slPackageResponse.Data[0]).To(Equal(clients.DataPackage{
				Id:   0,
				Name: "name0"}))
			Expect(slPackageResponse.Data[1]).To(Equal(clients.DataPackage{
				Id:   1,
				Name: "name1"}))
		})

		It("fails when BMP server /sl/packages fails", func() {
			fakeHttpClient.DoRawHttpRequestError = errors.New("fake-error")

			_, err := bmpClient.SlPackages()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("#stemcells", func() {
		BeforeEach(func() {
			fakeHttpClient.DoRawHttpRequestResponse, err = common.ReadJsonTestFixtures("..", "bmp", "Stemcells.json")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an array of stemcells", func() {
			stemcellsResponse, err := bmpClient.Stemcells()
			Expect(err).ToNot(HaveOccurred())

			Expect(len(stemcellsResponse.Stemcell)).To(Equal(2))
			Expect(stemcellsResponse.Stemcell[0]).To(Equal(
				"fake-stemcell-0"))
			Expect(stemcellsResponse.Stemcell[1]).To(Equal(
				"fake-stemcell-1"))
		})

		It("fails when BMP server /stemcells fails", func() {
			fakeHttpClient.DoRawHttpRequestError = errors.New("fake-error")

			_, err := bmpClient.Stemcells()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("#SlPackageOptions", func() {
		BeforeEach(func() {
			fakeHttpClient.DoRawHttpRequestResponse, err = common.ReadJsonTestFixtures("..", "bmp", "SlPackageOptions.json")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns DataPackageOptions info", func() {
			slPackageOptionsResponse, err := bmpClient.SlPackageOptions("fake-id")
			Expect(err).ToNot(HaveOccurred())

			Expect(len(slPackageOptionsResponse.Data.Category)).To(Equal(2))

			Options1 := []clients.Option{
				clients.Option{Id: 0, Description: "description0"},
				clients.Option{Id: 1, Description: "description1"},
			}

			Expect(slPackageOptionsResponse.Data.Category[0]).To(Equal(clients.Category{
				Code:     "code0",
				Name:     "name0",
				Options:  Options1,
				Required: true}))

			Options2 := []clients.Option{
				clients.Option{Id: 0, Description: "description0"},
			}

			Expect(slPackageOptionsResponse.Data.Category[1]).To(Equal(clients.Category{
				Code:     "code1",
				Name:     "name1",
				Options:  Options2,
				Required: false}))

			Expect(len(slPackageOptionsResponse.Data.Datacenter)).To(Equal(2))
			Expect(slPackageOptionsResponse.Data.Datacenter[0]).To(Equal(
				"datacenter0 - location0"))
			Expect(slPackageOptionsResponse.Data.Datacenter[1]).To(Equal(
				"datacenter1 - location1"))
		})

		It("fails when BMP server /sl/package/{packageid}/options fails", func() {
			fakeHttpClient.DoRawHttpRequestError = errors.New("fake-error")

			_, err := bmpClient.SlPackageOptions("fake-id")
			Expect(err).To(HaveOccurred())
		})
	})
})
