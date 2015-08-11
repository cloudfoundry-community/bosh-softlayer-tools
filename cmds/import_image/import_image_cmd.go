package import_image

import (
	"errors"
	"fmt"

	common "github.com/maximilien/bosh-softlayer-stemcells/common"

	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

type ImportImageCmd struct {
	options common.Options
	client  softlayer.Client

	Name      string
	Note      string
	OsRefCode string
	Uri       string

	Id   int
	Uuid string
}

func NewImportImageCmd(options common.Options, client softlayer.Client) (*ImportImageCmd, error) {
	return &ImportImageCmd{
		options: options,
		client:  client,

		Name:      options.NameFlag,
		Note:      options.NoteFlag,
		OsRefCode: options.OsRefCodeFlag,
		Uri:       options.UriFlag,
	}, nil
}

func (cmd *ImportImageCmd) Println(a ...interface{}) (int, error) {
	fmt.Println(a)

	return 0, nil
}

func (cmd *ImportImageCmd) Printf(msg string, a ...interface{}) (int, error) {
	fmt.Printf(msg, a)

	return 0, nil
}

func (cmd *ImportImageCmd) Options() common.Options {
	return cmd.options
}

func (cmd *ImportImageCmd) CheckOptions() error {
	if cmd.OsRefCode == "" {
		return errors.New("stemcells: must pass an OS ref code")
	}

	if cmd.Uri == "" {
		return errors.New("stemcells: must pass a URI")
	}

	return nil
}

func (cmd *ImportImageCmd) Run() error {
	vgbdtgService, err := cmd.client.GetSoftLayer_Virtual_Guest_Block_Device_Template_Group_Service()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get virtual guest block device template group service from softlayer-go service: `%s`", err.Error()))
	}

	configuration := sldatatypes.SoftLayer_Container_Virtual_Guest_Block_Device_Template_Configuration{
		Name: cmd.Name,
		Note: cmd.Note,
		OperatingSystemReferenceCode: cmd.OsRefCode,
		Uri: cmd.Uri,
	}

	vgbdtgObject, err := vgbdtgService.CreateFromExternalSource(configuration)
	if err != nil {
		return errors.New(fmt.Sprintf("Problem creating image template from external source: `%s`", err.Error()))
	}

	cmd.Id = vgbdtgObject.Id
	cmd.Uuid = vgbdtgObject.GlobalIdentifier

	return nil
}
