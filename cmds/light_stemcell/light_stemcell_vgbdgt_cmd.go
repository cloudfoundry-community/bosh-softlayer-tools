package light_stemcell

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cloudfoundry-community/bosh-softlayer-tools/common"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	"github.com/maximilien/softlayer-go/softlayer"
)

type LightStemcellVGBDTGCmd struct {
	options common.Options

	path    string
	name    string
	version string

	stemcellInfoFilename string

	lightStemcellInfo LightStemcellInfo

	infrastructure  string
	hypervisor      string
	osName          string
	stemcellFormats []string

	client softlayer.Client
}

func NewLightStemcellVGBDGTCmd(options common.Options, client softlayer.Client) *LightStemcellVGBDTGCmd {
	stemcellFormats := strings.Split(options.StemcellFormatsFlag, ",")

	cmd := &LightStemcellVGBDTGCmd{
		options: options,

		path:    options.LightStemcellPathFlag,
		name:    options.NameFlag,
		version: options.VersionFlag,

		stemcellInfoFilename: options.StemcellInfoFilenameFlag,

		infrastructure:  options.InfrastructureFlag,
		hypervisor:      options.HypervisorFlag,
		osName:          options.OsNameFlag,
		stemcellFormats: stemcellFormats,

		client: client,
	}

	cmd.updateLightStemcellInfo()

	return cmd
}

func (cmd *LightStemcellVGBDTGCmd) Println(a ...interface{}) (int, error) {
	fmt.Println(a)

	return 0, nil
}

func (cmd *LightStemcellVGBDTGCmd) Printf(msg string, a ...interface{}) (int, error) {
	fmt.Printf(msg, a)

	return 0, nil
}

func (cmd *LightStemcellVGBDTGCmd) Options() common.Options {
	return cmd.options
}

func (cmd *LightStemcellVGBDTGCmd) CheckOptions() error {
	if cmd.version == "" {
		return errors.New("light stemcell: must pass a version")
	}

	if cmd.stemcellInfoFilename == "" {
		return errors.New("light stemcell: must pass a path to stemcell-info.json")
	}

	return nil
}

func (cmd *LightStemcellVGBDTGCmd) Run() error {
	cmd.updateLightStemcellInfo()

	softLayerStemcellInfo, err := cmd.createSoftLayerStemcellInfo()
	if err != nil {
		return err
	}

	cmd.path, err = cmd.Create(softLayerStemcellInfo.Id)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *LightStemcellVGBDTGCmd) GetStemcellPath() string {
	return cmd.path
}

func (cmd *LightStemcellVGBDTGCmd) GetLightStemcellInfo() LightStemcellInfo {
	return cmd.lightStemcellInfo
}

func (cmd *LightStemcellVGBDTGCmd) Create(vgbdtgId int) (string, error) {
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

func (cmd *LightStemcellVGBDTGCmd) updateLightStemcellInfo() {
	cmd.lightStemcellInfo.Infrastructure = cmd.infrastructure
	cmd.lightStemcellInfo.Architecture = "x86_64"
	cmd.lightStemcellInfo.RootDeviceName = "/dev/xvda"

	cmd.lightStemcellInfo.Version = cmd.version
	cmd.lightStemcellInfo.Hypervisor = cmd.hypervisor
	cmd.lightStemcellInfo.OsName = cmd.osName

	// separating items-with comma
	cmd.lightStemcellInfo.StemcellFormats = cmd.stemcellFormats
}

func (cmd *LightStemcellVGBDTGCmd) createSoftLayerStemcellInfo() (SoftLayerStemcellInfo, error) {
	var softLayerStemcellInfo SoftLayerStemcellInfo

	slInfoFile, err := ioutil.ReadFile(cmd.stemcellInfoFilename)
	if err != nil {
		return softLayerStemcellInfo, errors.New(fmt.Sprintf("Could not read from SoftLayer info file: `%s`", err.Error()))
	}

	json.Unmarshal(slInfoFile, &softLayerStemcellInfo)

	return softLayerStemcellInfo, nil
}

func (cmd *LightStemcellVGBDTGCmd) findInVirtualGuestBlockDeviceTemplateGroups(vgdtgId int) (sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group, bool, error) {
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

func (cmd *LightStemcellVGBDTGCmd) buildLightStemcellWithVirtualGuestBlockDeviceTemplateGroup(vgdtgGroup sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group) (string, error) {
	datacenterName, err := cmd.findDatacenterFromVirtualGuestBlockDeviceTemplateGroup(vgdtgGroup)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Getting datacenter name from softlayer-go client: `%s`", err.Error()))
	}

	lightStemcellMF := LightStemcellMF{
		Name:            GenerateStemcellName(cmd.lightStemcellInfo),
		Version:         cmd.lightStemcellInfo.Version,
		BoshProtocol:    1, //Must be defaulted to 1 for legacy reasons (no other values supported)
		Sha1:            base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(fmt.Sprintf("%d:%s", vgdtgGroup.Id, vgdtgGroup.GlobalIdentifier)))),
		OperatingSystem: cmd.lightStemcellInfo.OsName,
		CloudProperties: CloudProperties{
			Infrastructure:       cmd.lightStemcellInfo.Infrastructure,
			Architecture:         cmd.lightStemcellInfo.Architecture,
			RootDeviceName:       cmd.lightStemcellInfo.RootDeviceName,
			VirtualDiskImageId:   vgdtgGroup.Id,
			VirtualDiskImageUuid: vgdtgGroup.GlobalIdentifier,
			DatacenterName:       datacenterName,
		},
		StemcellFormats: cmd.lightStemcellInfo.StemcellFormats,
	}

	return GenerateLightStemcellTarball(lightStemcellMF, cmd.lightStemcellInfo, cmd.path)
}

func (cmd *LightStemcellVGBDTGCmd) findDatacenterFromVirtualGuestBlockDeviceTemplateGroup(vgdtgGroup sldatatypes.SoftLayer_Virtual_Guest_Block_Device_Template_Group) (string, error) {
	vgdtgGroupService, err := cmd.client.GetSoftLayer_Virtual_Guest_Block_Device_Template_Group_Service()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get find data center for virtual guest device block template group '%d', got error from softlayer-go client: `%s`", vgdtgGroup.Id, err.Error()))
	}

	locations, err := vgdtgGroupService.GetDatacenters(vgdtgGroup.Id)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Could not get datacenters for virtual guest device block template group '%d' from softlayer-go client: `%s`", vgdtgGroup.Id, err.Error()))
	}

	if len(locations) > 0 {
		return locations[0].Name, nil
	}

	return "", nil
}
