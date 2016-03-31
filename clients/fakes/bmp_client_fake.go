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
