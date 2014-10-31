package stemcells

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

type lightStemcellVGBDTGCmd struct {
	ligthStemcellsPath string
	lightStemcellInfo  LightStemcellInfo
	client             softlayer.Client
}

func NewLightStemcellVGBDGTCmd(stemcellsPath string, lightStemcellInfo LightStemcellInfo, client softlayer.Client) *lightStemcellVGBDTGCmd {
	return &lightStemcellVGBDTGCmd{
		ligthStemcellsPath: stemcellsPath,
		lightStemcellInfo:  lightStemcellInfo,
		client:             client,
	}
}

func (cmd *lightStemcellVGBDTGCmd) GetStemcellsPath() string {
	return cmd.ligthStemcellsPath
}

func (cmd *lightStemcellVGBDTGCmd) GetLightStemcellInfo() LightStemcellInfo {
	return cmd.lightStemcellInfo
}

func (cmd *lightStemcellVGBDTGCmd) Create(vgbdtgId int) (string, error) {
	vgbdtGroup, found, err := cmd.findInVirtualGuestBlockDeviceTemplateGroups(vgbdtgId)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get virtual guest block device template group '%d' from softlayer-go service: `%s`", vgbdtgId, err.Error()))
	}

	if found == true {
		vgbdtGroupService, err := cmd.client.GetSoftLayer_Virtual_Guest_Block_Device_Template_Group_Service()
		if err != nil {
			return "", errors.New(fmt.Sprintf("Could not get virtual guest block device template group service from softlayer-go service: `%s`", err.Error()))
		}

		object, err := vgbdtGroupService.GetObject(vgbdtGroup.Id)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Could not get virtual guest block device template group object '%d' got softlayer-go service: `%s`", vgbdtgId, err.Error()))
		}

		return cmd.buildLightStemcellWithVirtualGuestBlockDeviceTemplateGroup(object)
	}

	return "", errors.New(fmt.Sprintf("Could not get virtual guest block device template group '%d'", vgbdtgId))
}

// Private methods

func (cmd *lightStemcellVGBDTGCmd) findInVirtualGuestBlockDeviceTemplateGroups(vgdtgId int) (sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, bool, error) {
	accountService, err := cmd.client.GetSoftLayer_Account_Service()
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

func (cmd *lightStemcellVGBDTGCmd) buildLightStemcellWithVirtualGuestBlockDeviceTemplateGroup(vgdtgGroup sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group) (string, error) {
	datacenterName, err := cmd.findDatacenterFromVirtualGuestBlockDeviceTemplateGroup(vgdtgGroup)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Getting datacenter name from softlayer-go client: `%s`", err.Error()))
	}

	lightStemcellMF := LightStemcellMF{
		Name:         GenerateStemcellName(cmd.lightStemcellInfo),
		Version:      cmd.lightStemcellInfo.Version,
		BoshProtocol: 1, //Must be defaulted to 1 for legacy reasons (no other values supported)
		Sha1:         base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(fmt.Sprintf("%d:%s", vgdtgGroup.Id, vgdtgGroup.GlobalIdentifier)))),
		CloudProperties: CloudProperties{
			Infrastructure:       cmd.lightStemcellInfo.Infrastructure,
			Architecture:         cmd.lightStemcellInfo.Architecture,
			RootDeviceName:       cmd.lightStemcellInfo.RootDeviceName,
			VirtualDiskImageId:   vgdtgGroup.Id,
			VirtualDiskImageUuid: vgdtgGroup.GlobalIdentifier,
			DatacenterName:       datacenterName,
		},
	}

	return GenerateLightStemcellTarball(lightStemcellMF, cmd.lightStemcellInfo, cmd.ligthStemcellsPath)
}

func (cmd *lightStemcellVGBDTGCmd) findDatacenterFromVirtualGuestBlockDeviceTemplateGroup(vgdtgGroup sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group) (string, error) {
	vgdtgGroupService, err := cmd.client.GetSoftLayer_Virtual_Guest_Block_Device_Template_Group_Service()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get find data center for virtual guest device block template group '%d', got error from softlayer-go client: `%s`", vgdtgGroup.Id, err.Error()))
	}

	locations, err := vgdtgGroupService.GetDatacenters(vgdtgGroup.Id)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not datacenters for virtual guest device block template group '%d' from softlayer-go client: `%s`", vgdtgGroup.Id, err.Error()))
	}

	if len(locations) > 0 {
		return locations[0].Name, nil
	}

	return "", nil
}
