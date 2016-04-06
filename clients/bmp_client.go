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
	username   string
	password   string
	url        string
	httpClient softlayer.HttpClient
}

func NewBmpClient(username, password, url string, hClient softlayer.HttpClient) *bmpClient {
	var httpClient softlayer.HttpClient
	if hClient == nil {
		httpClient = slclient.NewHttpClient(username, password, url, "")
	} else {
		httpClient = hClient
	}

	return &bmpClient{
		username:   username,
		password:   password,
		url:        url,
		httpClient: httpClient,
	}
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
