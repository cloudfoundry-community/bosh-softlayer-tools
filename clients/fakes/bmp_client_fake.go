package fakes

import (
	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
)

type FakeBmpClient struct {
	Username string
	Password string
	Url      string

	InfoResponse clients.InfoResponse
	InfoErr      error

	SlPackagesResponse clients.SlPackagesResponse
	SlPackagesErr      error

	StemcellResponse clients.StemcellsResponse
	StemcellErr      error

	SlPackageOptionsResponse clients.SlPackageOptionsResponse
	SlPackageOptionsErr      error
}

func NewFakeBmpClient(username, password, url string) *FakeBmpClient {
	return &FakeBmpClient{
		Username: username,
		Password: password,
		Url:      url,
	}
}

func (bc *FakeBmpClient) Info() (clients.InfoResponse, error) {
	return bc.InfoResponse, bc.InfoErr
}

func (bc *FakeBmpClient) SlPackages() (clients.SlPackagesResponse, error) {
	return bc.SlPackagesResponse, bc.SlPackagesErr
}

func (bc *FakeBmpClient) Stemcells() (clients.StemcellsResponse, error) {
	return bc.StemcellResponse, bc.StemcellErr
}

func (bc *FakeBmpClient) SlPackageOptions(packageId string) (clients.SlPackageOptionsResponse, error) {
	return bc.SlPackageOptionsResponse, bc.SlPackageOptionsErr
}
