package stemcells

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

type lightStemcellVDICmd struct {
	lightStemcellsPath string
	lightStemcellInfo  LightStemcellInfo
	client             softlayer.Client
}

func NewLightStemcellVDICmd(stemcellsPath string, lightStemcellInfo LightStemcellInfo, client softlayer.Client) *lightStemcellVDICmd {
	return &lightStemcellVDICmd{
		lightStemcellsPath: stemcellsPath,
		lightStemcellInfo:  lightStemcellInfo,
		client:             client,
	}
}

func (cmd *lightStemcellVDICmd) GetStemcellsPath() string {
	return cmd.lightStemcellsPath
}

func (cmd *lightStemcellVDICmd) GetLightStemcellInfo() LightStemcellInfo {
	return cmd.lightStemcellInfo
}

func (cmd *lightStemcellVDICmd) Create(vdImageId int) (string, error) {
	virtualDiskImageService, err := cmd.client.GetSoftLayer_Virtual_Disk_Image_Service()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get SoftLayer_Virtual_Disk_Image_Service from softlayer-go client: `%s`", err.Error()))
	}

	virtualDiskImage, err := virtualDiskImageService.GetObject(vdImageId)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get Virtual_Disk_Image from softlayer-go service: `%s`", err.Error()))
	}

	emptyVirtualDiskImage := sldatatypes.SoftLayer_Virtual_Disk_Image{}
	if virtualDiskImage == emptyVirtualDiskImage {
		vdImage, found, err := cmd.findInVirtualDiskImages(vdImageId)
		if err != nil {
			return "", err
		}

		if found == false {
			return "", errors.New(fmt.Sprintf("Did not find SoftLayer virtual disk image with ID '%d'", vdImageId))
		}

		virtualDiskImage = vdImage
	}

	return cmd.buildLightStemcellWithVirtualDiskImage(virtualDiskImage)
}

// Private methods

func (cmd *lightStemcellVDICmd) findInVirtualDiskImages(vdImageId int) (sldatatypes.SoftLayer_Virtual_Disk_Image, bool, error) {
	accountService, err := cmd.client.GetSoftLayer_Account_Service()
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

func (cmd *lightStemcellVDICmd) buildLightStemcellWithVirtualDiskImage(virtualDiskImage sldatatypes.SoftLayer_Virtual_Disk_Image) (string, error) {
	datacenterName, err := cmd.findDatacenterFromVirtualDiskImage(virtualDiskImage)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Getting datacenter name from softlayer-go client: `%s`", err.Error()))
	}

	lightStemcellMF := LightStemcellMF{
		Name:         GenerateStemcellName(cmd.lightStemcellInfo),
		Version:      cmd.lightStemcellInfo.Version,
		BoshProtocol: 1, //Must be defaulted to 1 for legacy reasons (no other values supported)
		Sha1:         base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(fmt.Sprintf("%d:%s", virtualDiskImage.Id, virtualDiskImage.Uuid)))),
		CloudProperties: CloudProperties{
			Infrastructure:       cmd.lightStemcellInfo.Infrastructure,
			Architecture:         cmd.lightStemcellInfo.Architecture,
			RootDeviceName:       cmd.lightStemcellInfo.RootDeviceName,
			VirtualDiskImageId:   virtualDiskImage.Id,
			VirtualDiskImageUuid: virtualDiskImage.Uuid,
			DatacenterName:       datacenterName,
		},
	}

	return GenerateLightStemcellTarball(lightStemcellMF, cmd.lightStemcellInfo, cmd.lightStemcellsPath)
}

func (cmd *lightStemcellVDICmd) findDatacenterFromVirtualDiskImage(virtualDiskImage sldatatypes.SoftLayer_Virtual_Disk_Image) (string, error) {
	accountService, err := cmd.client.GetSoftLayer_Account_Service()
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
