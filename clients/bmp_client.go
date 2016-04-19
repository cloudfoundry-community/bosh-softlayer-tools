package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	slclient "github.com/maximilien/softlayer-go/client"
	slcommon "github.com/maximilien/softlayer-go/common"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

type bmpClient struct {
	username string
	password string
	url      string

	configPath string

	httpClient softlayer.HttpClient
}

func NewBmpClient(username, password, url string, hClient softlayer.HttpClient, configPath string) *bmpClient {
	var httpClient softlayer.HttpClient
	if hClient == nil {
		httpClient = slclient.NewHttpClient(username, password, url, "")
	} else {
		httpClient = hClient
	}

	return &bmpClient{
		username: username,
		password: password,
		url:      url,

		configPath: configPath,

		httpClient: httpClient,
	}
}

func (bc *bmpClient) ConfigPath() string {
	return bc.configPath
}

func (bc *bmpClient) Info() (InfoResponse, error) {
	path := fmt.Sprintf("%s/%s", bc.url, "info")
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "GET", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /info on BMP server, error message '%s'", err.Error())
		return InfoResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /info on BMP server, HTTP error code: '%d'", errorCode)
		return InfoResponse{}, errors.New(errorMessage)
	}

	response := InfoResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return InfoResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) Bms(deploymentName string) (BmsResponse, error) {
	path := fmt.Sprintf("%s/%s/%s", bc.url, "/bms/", deploymentName)
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "GET", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /bms/'%s' on BMP server, error message '%s'", deploymentName, err.Error())
		return BmsResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /bms/'%s' on BMP server, HTTP error code: '%d'", deploymentName, errorCode)
		return BmsResponse{}, errors.New(errorMessage)
	}

	response := BmsResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return BmsResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) SlPackages() (SlPackagesResponse, error) {
	path := fmt.Sprintf("%s/%s", bc.url, "sl/packages")
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "GET", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /sl/packages on BMP server, error message '%s'", err.Error())
		return SlPackagesResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /info on BMP server, HTTP error code: '%d'", errorCode)
		return SlPackagesResponse{}, errors.New(errorMessage)
	}

	response := SlPackagesResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return SlPackagesResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) Stemcells() (StemcellsResponse, error) {
	path := fmt.Sprintf("%s/%s", bc.url, "/stemcells")
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "GET", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /stemcells on BMP server, error message '%s'", err.Error())
		return StemcellsResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /stemcells on BMP server, HTTP error code: '%d'", errorCode)
		return StemcellsResponse{}, errors.New(errorMessage)
	}

	response := StemcellsResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return StemcellsResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) SlPackageOptions(packageId string) (SlPackageOptionsResponse, error) {
	path := fmt.Sprintf("%s/sl/packages/%s/options", bc.url, packageId)
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "GET", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /sl/packages/'%s'/options on BMP server, error message '%s'", packageId, err.Error())
		return SlPackageOptionsResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /sl/packages/'%s'/options on BMP server, HTTP error code: '%d'", packageId, errorCode)
		return SlPackageOptionsResponse{}, errors.New(errorMessage)
	}

	response := SlPackageOptionsResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return SlPackageOptionsResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) Tasks(latest int) (TasksResponse, error) {
	path := fmt.Sprintf("%s/%s%d", bc.url, "/tasks?latest=", latest)
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "GET", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /tasks?latest='%d'on BMP server, error message '%s'", latest, err.Error())
		return TasksResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /tasks?latest='%d' on BMP server, HTTP error code: '%d'", latest, errorCode)
		return TasksResponse{}, errors.New(errorMessage)
	}

	response := TasksResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return TasksResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) TaskOutput(taskId int, level string) (TaskOutputResponse, error) {
	path := fmt.Sprintf("%s/task/%d/txt/%s", bc.url, taskId, level)
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "GET", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /task/'%d'/txt/'%s'on BMP server, error message '%s'", taskId, level, err.Error())
		return TaskOutputResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /task/'%d'/txt/'%s' on BMP server, HTTP error code: '%d'", taskId, level, errorCode)
		return TaskOutputResponse{}, errors.New(errorMessage)
	}

	response := TaskOutputResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return TaskOutputResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) UpdateStatus(serverId string, status string) (UpdateStatusResponse, error) {
	path := fmt.Sprintf("%s/baremetal/%s/%s", bc.url, serverId, status)
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "PUT", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /baremetal/'%s'/'%s' on BMP server, error message '%s'", serverId, status, err.Error())
		return UpdateStatusResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /baremetal/'%s'/'%s' on BMP server, HTTP error code: '%d'", serverId, status, errorCode)
		return UpdateStatusResponse{}, errors.New(errorMessage)
	}

	response := UpdateStatusResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return UpdateStatusResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) Login(username string, password string) (LoginResponse, error) {
	path := fmt.Sprintf("%s/login/%s/%s", bc.url, username, password)
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "GET", &bytes.Buffer{})
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /login/'%s'/'%s' on BMP server, error message '%s'", username, password, err.Error())
		return LoginResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /login/'%s'/'%s' on BMP server, HTTP error code: '%d'", username, password, errorCode)
		return LoginResponse{}, errors.New(errorMessage)
	}

	response := LoginResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return LoginResponse{}, errors.New(errorMessage)
	}

	return response, nil
}

func (bc *bmpClient) CreateBaremetal(createBaremetalInfo CreateBaremetalInfo) (CreateBaremetalResponse, error) {
	createBaremetalParameters := CreateBaremetalParameters{
		Parameters: createBaremetalInfo,
	}

	requestBody, err := json.Marshal(createBaremetalParameters)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to encode Json File, error message '%s'", err.Error())
		return CreateBaremetalResponse{}, errors.New(errorMessage)
	}

	path := fmt.Sprintf("%s/%s", bc.url, "sl/packages")
	responseBytes, errorCode, err := bc.httpClient.DoRawHttpRequest(path, "POST", bytes.NewBuffer(requestBody))
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not calls /baremetals on BMP server, error message '%s'", err.Error())
		return CreateBaremetalResponse{}, errors.New(errorMessage)
	}

	if slcommon.IsHttpErrorCode(errorCode) {
		errorMessage := fmt.Sprintf("bmp: could not call /baremetals on BMP server, HTTP error code: '%d'", errorCode)
		return CreateBaremetalResponse{}, errors.New(errorMessage)
	}

	response := CreateBaremetalResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode JSON response, err message '%s'", err.Error())
		return CreateBaremetalResponse{}, errors.New(errorMessage)
	}

	return response, nil
}
