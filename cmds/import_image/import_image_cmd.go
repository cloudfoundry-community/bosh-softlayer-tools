package import_image

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-community/bosh-softlayer-tools/common"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshretry "github.com/cloudfoundry/bosh-utils/retrystrategy"
	sldatatypes "github.com/maximilien/softlayer-go/data_types"
	"github.com/maximilien/softlayer-go/softlayer"
	"github.com/pivotal-golang/clock"
)

var (
	locations = []sldatatypes.SoftLayer_Location{
		{
			358694,
			"London 2",
			"lon02",
		},
	}
)

type ImportImageCmd struct {
	options common.Options
	client  softlayer.Client

	Public     bool
	Name       string
	Note       string
	PublicName string
	PublicNote string
	OsRefCode  string
	Uri        string

	Id   int
	Uuid string
}

func NewImportImageCmd(options common.Options, client softlayer.Client) (*ImportImageCmd, error) {
	return &ImportImageCmd{
		options: options,
		client:  client,

		Public:     options.PublicFlag,
		Name:       options.NameFlag,
		Note:       options.NoteFlag,
		PublicName: options.PublicNameFlag,
		PublicNote: options.PublicNoteFlag,
		OsRefCode:  options.OsRefCodeFlag,
		Uri:        options.UriFlag,
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

	if cmd.Public {
		execStmtRetryable := boshretry.NewRetryable(
			func() (bool, error) {
				id, err := vgbdtgService.CreatePublicArchiveTransaction(cmd.Id, cmd.PublicName, cmd.PublicNote, cmd.PublicNote, locations)
				if err != nil {
					return true, errors.New(fmt.Sprintf("There would be an active transaction in progress."))
				}

				cmd.Id = id
				cmd.Uuid = ""

				return false, nil
			})
		timeService := clock.NewClock()
		timeoutRetryStrategy := boshretry.NewTimeoutRetryStrategy(common.TIMEOUT, common.POLLING_INTERVAL, execStmtRetryable, timeService, boshlog.NewLogger(boshlog.LevelInfo))
		err = timeoutRetryStrategy.Try()
		if err != nil {
			return errors.New(fmt.Sprintf("Problem occurred when making image template public: `%s`", err.Error()))
		} else {
			return nil
		}
	}

	return nil
}
