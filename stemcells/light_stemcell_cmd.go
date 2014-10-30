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
	ligthStemcellPath, err := lsd.createLightStemcellFromVirtualDiskImage(virtualDiskImageId)
	if err != nil {
		return "", err
	}

	if ligthStemcellPath != "" {
		return ligthStemcellPath, nil
	}

	ligthStemcellPath, err = lsd.createLightStemcellFromVirtualGuestDeviceTemplateGroup(virtualDiskImageId)
	if err != nil {
		return "", err
	}

	if ligthStemcellPath != "" {
		return ligthStemcellPath, nil
	}

	return "", errors.New(fmt.Sprintf("Could not create SoftLayer light stemcell from ID: '%d'", virtualDiskImageId))
}

func (lsd *lightStemcellCmd) createLightStemcellFromVirtualDiskImage(vdImageId int) (string, error) {
	virtualDiskImageService, err := lsd.client.GetSoftLayer_Virtual_Disk_Image_Service()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get SoftLayer_Virtual_Disk_Image_Service from softlayer-go client: `%s`", err.Error()))
	}

	virtualDiskImage, err := virtualDiskImageService.GetObject(vdImageId)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get Virtual_Disk_Image from softlayer-go service: `%s`", err.Error()))
	}

	emptyVirtualDiskImage := sldatatypes.SoftLayer_Virtual_Disk_Image{}
	if virtualDiskImage == emptyVirtualDiskImage {
		vdImage, found, err := lsd.findInVirtualDiskImages(vdImageId)
		if err != nil {
			return "", err
		}

		if found == false {
			return "", errors.New(fmt.Sprintf("Did not find SoftLayer virtual disk image with ID '%d'", vdImageId))
		}

		virtualDiskImage = vdImage
	}

	return lsd.buildLightStemcellWithVirtualDiskImage(virtualDiskImage)
}

func (lsd *lightStemcellCmd) findInVirtualDiskImages(vdImageId int) (sldatatypes.SoftLayer_Virtual_Disk_Image, bool, error) {
	accountService, err := lsd.client.GetSoftLayer_Account_Service()
	if err != nil {
		return sldatatypes.SoftLayer_Virtual_Disk_Image{}, false, errors.New(fmt.Sprintf("Could not get SoftLayer_Account_Service from softlayer-go client: `%s`", err.Error()))
	}

	virtualDiskImages, err := accountService.GetVirtualDiskImages()
	if err != nil {
		return sldatatypes.SoftLayer_Virtual_Disk_Image{}, false, errors.New(fmt.Sprintf("Getting virtual disk images from softlayer-go client: `%s`", err.Error()))
	}

	for _, vdImage := range virtualDiskImages {
		if vdImage.Id == vdImageId {
			return vdImage, true, nil
		}
	}

	return sldatatypes.SoftLayer_Virtual_Disk_Image{}, false, nil
}

func (lsd *lightStemcellCmd) buildLightStemcellWithVirtualDiskImage(virtualDiskImage sldatatypes.SoftLayer_Virtual_Disk_Image) (string, error) {
	return "", nil
}

func (lsd *lightStemcellCmd) createLightStemcellFromVirtualGuestDeviceTemplateGroup(vgdtgId int) (string, error) {
	return "", nil
}

func (lsd *lightStemcellCmd) findInVirtualGuestDeviceTemplateGroups(vgdtgId int) (sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, bool, error) {
	accountService, err := lsd.client.GetSoftLayer_Account_Service()
	if err != nil {
		return sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group{}, false, errors.New(fmt.Sprintf("Could not get SoftLayer_Account_Service from softlayer-go client: `%s`", err.Error()))
	}

	vgdtgGroups, err := accountService.GetBlockDeviceTemplateGroups()
	if err != nil {
		return sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group{}, false, errors.New(fmt.Sprintf("Getting virtual guest device template groups from softlayer-go service: '%s'", err.Error()))
	}

	for _, vgdtgGroup := range vgdtgGroups {
		if vgdtgGroup.Id == vgdtgId {
			return vgdtgGroup, true, nil
		}
	}

	return sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group{}, false, nil
}

func (lsd *lightStemcellCmd) buildLightStemcellWithVirtualGuestDeviceTemplateGroup(vgdtgGroup sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group) (string, error) {
	return "", nil
}
