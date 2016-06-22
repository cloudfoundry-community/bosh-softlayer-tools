package light_stemcell

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

func GenerateStemcellName(info LightStemcellInfo) string {
	return fmt.Sprintf("bosh-%s-%s-%s-go_agent",
		info.Infrastructure,
		info.Hypervisor,
		info.OsName)
}

func GenerateStemcellFilelName(info LightStemcellInfo) string {
	return fmt.Sprintf("light-bosh-stemcell-%s-%s-%s-%s-go_agent",
		info.Version,
		info.Infrastructure,
		info.Hypervisor,
		info.OsName)
}

func GenerateLightStemcellTarball(lightStemcellMF LightStemcellMF, lightStemcellInfo LightStemcellInfo, lightStemcellsPath string) (string, error) {
	lightStemcellMFBytes, err := GenerateManifestMFBytesYAML(lightStemcellMF)

	lightStemcellMFFilePath := filepath.Join(lightStemcellsPath, "stemcell.MF")
	err = ioutil.WriteFile(lightStemcellMFFilePath, lightStemcellMFBytes, 0666)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not create stemcell.MF file, error: `%s`", err.Error()))
	}

	imageFilePath := filepath.Join(lightStemcellsPath, "image")
	err = ioutil.WriteFile(imageFilePath, []byte{}, 0666)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not create image file, error: `%s`", err.Error()))
	}

	lightStemcellFilePath := filepath.Join(lightStemcellsPath, fmt.Sprintf("%s.tgz", GenerateStemcellFilelName(lightStemcellInfo)))
	err = common.CreateTarball(lightStemcellFilePath, []string{lightStemcellMFFilePath, imageFilePath})
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not create tarball file, error: `%s`", err.Error()))
	}

	err = os.Remove(lightStemcellMFFilePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could clean up temp file '%s', error: `%s`", lightStemcellMFFilePath, err.Error()))
	}

	err = os.Remove(imageFilePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could clean up temp file '%s', error: `%s`", imageFilePath, err.Error()))
	}

	return lightStemcellFilePath, nil
}

func GenerateManifestMFBytesJSON(lightStemcellMF LightStemcellMF) ([]byte, error) {
	bytes, err := json.Marshal(&lightStemcellMF)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Could not marshall stemcell manifest data into JSON, error: `%s`", err.Error()))
	}

	return bytes, err
}

func GenerateManifestMFBytesYAML(lightStemcellMF LightStemcellMF) ([]byte, error) {
	bytes, err := yaml.Marshal(&lightStemcellMF)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Could not marshall stemcell manifest data into YML, error: `%s`", err.Error()))
	}

	return bytes, err
}
