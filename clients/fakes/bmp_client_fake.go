package fakes

import (
	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
)

type FakeBmpClient struct {
	Username string
	Password string
	Url      string

	ConfigPathString string

	InfoResponse clients.InfoResponse
	InfoErr      error

	BmsResponse clients.BmsResponse
	BmsErr      error

	SlPackagesResponse clients.SlPackagesResponse
	SlPackagesErr      error

	StemcellResponse clients.StemcellsResponse
	StemcellErr      error

	SlPackageOptionsResponse clients.SlPackageOptionsResponse
	SlPackageOptionsErr      error

	TasksResponse clients.TasksResponse
	TasksErr      error

	TaskOutputResponse clients.TaskOutputResponse
	TaskOutputErr      error

	UpdateStateResponse clients.UpdateStateResponse
	UpdateStateErr      error

	LoginResponse clients.LoginResponse
	LoginErr      error

	CreateBaremetalsResponse clients.CreateBaremetalsResponse
	CreateBaremetalsErr      error
}

func NewFakeBmpClient(username, password, url string, configPath string) *FakeBmpClient {
	return &FakeBmpClient{
		Username:         username,
		Password:         password,
		Url:              url,
		ConfigPathString: configPath,
	}
}

func (bc *FakeBmpClient) ConfigPath() string {
	return bc.ConfigPathString
}

func (bc *FakeBmpClient) Info() (clients.InfoResponse, error) {
	return bc.InfoResponse, bc.InfoErr
}

func (bc *FakeBmpClient) Bms(deploymentName string) (clients.BmsResponse, error) {
	return bc.BmsResponse, bc.BmsErr
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

func (bc *FakeBmpClient) Tasks(latest uint) (clients.TasksResponse, error) {
	return bc.TasksResponse, bc.TasksErr
}

func (bc *FakeBmpClient) TaskOutput(taskID uint, level string) (clients.TaskOutputResponse, error) {
	return bc.TaskOutputResponse, bc.TaskOutputErr
}

func (bc *FakeBmpClient) UpdateState(serverId string, status string) (clients.UpdateStateResponse, error) {
	return bc.UpdateStateResponse, bc.UpdateStateErr
}

func (bc *FakeBmpClient) Login(username string, password string) (clients.LoginResponse, error) {
	return bc.LoginResponse, bc.LoginErr
}

func (bc *FakeBmpClient) CreateBaremetals(CreateBaremetalsInfo clients.CreateBaremetalsInfo, DryRun bool) (clients.CreateBaremetalsResponse, error) {
	return bc.CreateBaremetalsResponse, bc.CreateBaremetalsErr
}
