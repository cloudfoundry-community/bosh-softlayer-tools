package light_stemcell

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

type LightStemcellVDICmd struct {
	options common.Options

	path    string
	name    string
	version string

	stemcellInfoFilename string

	lightStemcellInfo LightStemcellInfo

	infrastructure string
	hypervisor     string
	osName         string

	stemcellFormats []string

	client softlayer.Client
}

func NewLightStemcellVDICmd(options common.Options, client softlayer.Client) *LightStemcellVDICmd {
	stemcellFormats := strings.Split(options.StemcellFormatsFlag, ",")

	cmd := &LightStemcellVDICmd{
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

func (cmd *LightStemcellVDICmd) Println(a ...interface{}) (int, error) {
	fmt.Println(a)

	return 0, nil
}

func (cmd *LightStemcellVDICmd) Printf(msg string, a ...interface{}) (int, error) {
	fmt.Printf(msg, a)

	return 0, nil
}

func (cmd *LightStemcellVDICmd) Options() common.Options {
	return cmd.options
}

func (cmd *LightStemcellVDICmd) CheckOptions() error {
	if cmd.version == "" {
		return errors.New("light stemcell: must pass a version")
	}

	if cmd.stemcellInfoFilename == "" {
		return errors.New("light stemcell: must pass a path to stemcell-info.json")
	}

	return nil
}

func (cmd *LightStemcellVDICmd) Run() error {
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

func (cmd *LightStemcellVDICmd) GetStemcellPath() string {
	return cmd.path
}

func (cmd *LightStemcellVDICmd) GetLightStemcellInfo() LightStemcellInfo {
	return cmd.lightStemcellInfo
}

func (cmd *LightStemcellVDICmd) Create(vdImageId int) (string, error) {
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

func (cmd *LightStemcellVDICmd) updateLightStemcellInfo() {
	cmd.lightStemcellInfo.Infrastructure = cmd.infrastructure
	cmd.lightStemcellInfo.Architecture = "x86_64"
	cmd.lightStemcellInfo.RootDeviceName = "/dev/xvda"

	cmd.lightStemcellInfo.Version = cmd.version
	cmd.lightStemcellInfo.Hypervisor = cmd.hypervisor
	cmd.lightStemcellInfo.OsName = cmd.osName
	cmd.lightStemcellInfo.StemcellFormats = cmd.stemcellFormats
}

func (cmd *LightStemcellVDICmd) createSoftLayerStemcellInfo() (SoftLayerStemcellInfo, error) {
	var softLayerStemcellInfo SoftLayerStemcellInfo

	slInfoFile, err := ioutil.ReadFile(cmd.stemcellInfoFilename)
	if err != nil {
		return softLayerStemcellInfo, errors.New(fmt.Sprintf("Could not read from SoftLayer info file: `%s`", err.Error()))
	}

	err = json.Unmarshal(slInfoFile, &softLayerStemcellInfo)
	if err != nil {
		return softLayerStemcellInfo, errors.New(fmt.Sprintf("Could not unmarshal softLayerStemcellInfo: `%s`", err.Error()))
	}

	return softLayerStemcellInfo, nil
}

func (cmd *LightStemcellVDICmd) findInVirtualDiskImages(vdImageId int) (sldatatypes.SoftLayer_Virtual_Disk_Image, bool, error) {
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

func (cmd *LightStemcellVDICmd) buildLightStemcellWithVirtualDiskImage(virtualDiskImage sldatatypes.SoftLayer_Virtual_Disk_Image) (string, error) {
	datacenterName, err := cmd.findDatacenterFromVirtualDiskImage(virtualDiskImage)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Getting datacenter name from softlayer-go client: `%s`", err.Error()))
	}

	lightStemcellMF := LightStemcellMF{
		Name:            GenerateStemcellName(cmd.lightStemcellInfo),
		Version:         cmd.lightStemcellInfo.Version,
		BoshProtocol:    1, //Must be defaulted to 1 for legacy reasons (no other values supported)
		Sha1:            base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(fmt.Sprintf("%d:%s", virtualDiskImage.Id, virtualDiskImage.Uuid)))),
		OperatingSystem: cmd.lightStemcellInfo.OsName,
		CloudProperties: CloudProperties{
			Infrastructure:       cmd.lightStemcellInfo.Infrastructure,
			Version:              cmd.lightStemcellInfo.Version,
			Architecture:         cmd.lightStemcellInfo.Architecture,
			RootDeviceName:       cmd.lightStemcellInfo.RootDeviceName,
			VirtualDiskImageId:   virtualDiskImage.Id,
			VirtualDiskImageUuid: virtualDiskImage.Uuid,
			DatacenterName:       datacenterName,
		},
		StemcellFormats: cmd.lightStemcellInfo.StemcellFormats,
	}

	return GenerateLightStemcellTarball(lightStemcellMF, cmd.lightStemcellInfo, cmd.path)
}

func (cmd *LightStemcellVDICmd) findDatacenterFromVirtualDiskImage(virtualDiskImage sldatatypes.SoftLayer_Virtual_Disk_Image) (string, error) {
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
