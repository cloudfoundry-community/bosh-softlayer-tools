package stemcells

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	common "github.com/maximilien/bosh-softlayer-stemcells/common"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"

	yaml "github.com/fraenkel/candiedyaml"
)

type lightStemcellCmd struct {
	ligthStemcellsPath string
	lightStemcellInfo  LightStemcellInfo
	client             softlayer.Client
}

type LightStemcellMF struct {
	Name            string          `json:"name" yaml:"name"`
	Version         string          `json:"version" yaml:"version"`
	BoshProtocol    int             `json:"bosh_protocol" yaml:"bosh_protocol"`
	Sha1            string          `json:"sha1" yaml:"sha1"`
	CloudProperties CloudProperties `json:"cloud_properties" yaml:"cloud_properties"`
}

type CloudProperties struct {
	Infrastructure string `json:"infrastructure" yaml:"infrastructure"`
	Architecture   string `json:"architecture" yaml:"architecture"`
	RootDeviceName string `json:"root_device_name" yaml:"root_device_name"`

	//SoftLayer-specific properties
	VirtualDiskImageId   int    `json:"virtual-disk-image-id" yaml:"virtual-disk-image-id"`
	VirtualDiskImageUuid string `json:"virtual-disk-image-uuid" yaml:"virtual-disk-image-uuid"`
	DatacenterName       string `json:"datacenter-name" yaml:"datacenter-name"`
}

func NewLightStemcellCmd(stemcellsPath string, lightStemcellInfo LightStemcellInfo, client softlayer.Client) *lightStemcellCmd {
	return &lightStemcellCmd{
		ligthStemcellsPath: stemcellsPath,
		lightStemcellInfo:  lightStemcellInfo,
		client:             client,
	}
}

func (lsd *lightStemcellCmd) GenerateStemcellName(info LightStemcellInfo) string {
	return fmt.Sprintf("light-bosh-stemcell-%s-%s-%s-%s-go_agent",
		lsd.lightStemcellInfo.Version,
		lsd.lightStemcellInfo.Infrastructure,
		lsd.lightStemcellInfo.Hypervisor,
		lsd.lightStemcellInfo.OsName)
}

func (lsd *lightStemcellCmd) GetStemcellsPath() string {
	return lsd.ligthStemcellsPath
}

func (lsd *lightStemcellCmd) GetLightStemcellInfo() LightStemcellInfo {
	return lsd.lightStemcellInfo
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

// Private methods

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
	datacenterName, err := lsd.findDatacenterFromVirtualDiskImage(virtualDiskImage)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Getting datacenter name from softlayer-go client: `%s`", err.Error()))
	}

	lightStemcellMF := LightStemcellMF{
		Name:         lsd.GenerateStemcellName(lsd.lightStemcellInfo),
		Version:      lsd.lightStemcellInfo.Version,
		BoshProtocol: 1, //Must be defaulted to 1 for legacy reasons (no other values supported)
		Sha1:         base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(fmt.Sprintf("%d:%s", virtualDiskImage.Id, virtualDiskImage.Uuid)))),
		CloudProperties: CloudProperties{
			Infrastructure:       lsd.lightStemcellInfo.Infrastructure,
			Architecture:         lsd.lightStemcellInfo.Architecture,
			RootDeviceName:       lsd.lightStemcellInfo.RootDeviceName,
			VirtualDiskImageId:   virtualDiskImage.Id,
			VirtualDiskImageUuid: virtualDiskImage.Uuid,
			DatacenterName:       datacenterName,
		},
	}

	lightStemcellMFBytes, err := lsd.generateManifestMFBytesYAML(lightStemcellMF)

	lightStemcellMFFilePath := filepath.Join(lsd.ligthStemcellsPath, "stemcell.MF")
	err = common.CreateFile(lightStemcellMFFilePath, lightStemcellMFBytes)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not create stemcell.MF file, error: `%s`", err.Error()))
	}

	imageFilePath := filepath.Join(lsd.ligthStemcellsPath, "image")
	err = common.CreateFile(imageFilePath, []byte{})
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not create image file, error: `%s`", err.Error()))
	}

	lightStemcellFilePath := filepath.Join(lsd.ligthStemcellsPath, fmt.Sprintf("%s.tgz", lsd.GenerateStemcellName(lsd.lightStemcellInfo)))
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

func (lsd *lightStemcellCmd) generateManifestMFBytesJSON(lightStemcellMF LightStemcellMF) ([]byte, error) {
	bytes, err := json.Marshal(&lightStemcellMF)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Could not marshall stemcell manifest data into JSON, error: `%s`", err.Error()))
	}

	return bytes, err
}

func (lsd *lightStemcellCmd) generateManifestMFBytesYAML(lightStemcellMF LightStemcellMF) ([]byte, error) {
	bytes, err := yaml.Marshal(&lightStemcellMF)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Could not marshall stemcell manifest data into YML, error: `%s`", err.Error()))
	}

	return bytes, err
}

func (lsd *lightStemcellCmd) findDatacenterFromVirtualDiskImage(virtualDiskImage sldatatypes.SoftLayer_Virtual_Disk_Image) (string, error) {
	accountService, err := lsd.client.GetSoftLayer_Account_Service()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get SoftLayer_Account_Service from softlayer-go client: `%s`", err.Error()))
	}

	locations, err := accountService.GetDatacentersWithSubnetAllocations()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get SoftLayer_Account_Service#GetDatacentersWithSubnetAllocations from softlayer-go client: `%s`", err.Error()))
	}

	if len(locations) > 0 {
		return locations[0].Name, nil
	}

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

func (lsd *lightStemcellCmd) findDatacenterFromVVirtualGuestDeviceTemplateGroup(vgdtgGroup sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group) (string, error) {
	return "", nil
}
