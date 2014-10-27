package stemcells

import (
	"errors"
	"fmt"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

type lightStemcellCmd struct {
	LightStemcellMF LightStemcellMF

	stemcellsPath string
	client        softlayer.Client
}

type LightStemcellMF struct {
	Name            string          `json:"name"`
	Version         string          `json:"name"`
	BoshProtocol    int             `json:"bosh_protocol"`
	Sha1            string          `json:"sha1"`
	CloudProperties CloudProperties `json:"cloud_properties"`
}

type CloudProperties struct {
	Infrastructure      string              `json:"infrastructure"`
	Architecture        string              `json:"architecture"`
	RootDeviceName      string              `json:"root_device_name"`
	SoftlayerProperties SoftlayerProperties `json:"softlayer"`
}

type SoftlayerProperties struct {
	VirtualDiskImageId   string `json:"virtual-disk-image-id"`
	VirtualDiskImageUuid string `json:"virtual-disk-image-uuid"`
	DatacenterName       string `json:"datacenter-name"`
}

func NewLightStemcellCmd(stemcellsPath string, client softlayer.Client) LightStemcellCmd {
	return &lightStemcellCmd{
		stemcellsPath: stemcellsPath,
		client:        client,
	}
}

func (lsd *lightStemcellCmd) GetLigthStemcellMF() LightStemcellMF {
	return lsd.LightStemcellMF
}

func (lsd *lightStemcellCmd) GetStemcellsPath() string {
	return lsd.stemcellsPath
}

func (lsd *lightStemcellCmd) Create(virtualDiskImageId int) (string, error) {
	virtualDiskImageService, err := lsd.client.GetSoftLayer_Virtual_Disk_Image_Service()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get SoftLayer_Virtual_Disk_Image_Service from softlayer-go client: `%s`", err.Error()))
	}

	virtualDiskImage, err := virtualDiskImageService.GetObject(virtualDiskImageId)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get Virtual_Disk_Image from softlayer-go service: `%s`", err.Error()))
	}

	ligthStemcellPath, err := lsd.createLightStemcell(virtualDiskImage)
	if err != nil {
		return "", err
	}

	return ligthStemcellPath, nil
}

func (lsd *lightStemcellCmd) createLightStemcell(virtualDiskImage sldatatypes.SoftLayer_Virtual_Disk_Image) (string, error) {
	return "", nil
}
