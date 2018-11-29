package bmp

import (
	"errors"

	"github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/common"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/common/fakes"
)

type provisioningBaremetalCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewProvisioningBaremetalCommand(options cmds.Options, bmpClient clients.BmpClient) provisioningBaremetalCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return provisioningBaremetalCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func NewFakeProvisioningBaremetalCommand(options cmds.Options, bmpClient clients.BmpClient, userInput string) provisioningBaremetalCommand {
	FakeUi := fakes.NewFakeUi()
	FakeUi.UserInput = userInput

	return provisioningBaremetalCommand{
		options:   options,
		ui:        FakeUi,
		printer:   common.NewDefaultPrinter(FakeUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd provisioningBaremetalCommand) Name() string {
	return "provisioning-baremetal"
}

func (cmd provisioningBaremetalCommand) Description() string {
	return `provisioning a baremetal with specific stemcell, netboot image`
}

func (cmd provisioningBaremetalCommand) Usage() string {
	return "bmp provisioning-baremetal --vmprefix <vm-prefix> --stemcell <bm-stemcell> --netbootimage <bm-netboot-image>"
}

func (cmd provisioningBaremetalCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd provisioningBaremetalCommand) Validate() (bool, error) {
	_, err := cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	if err != nil {
		return false, err
	}

	if cmd.options.NetbootImage == "" {
		return false, errors.New("cannot have empty netboot image")
	}

	if cmd.options.VMPrefix == "" {
		return false, errors.New("cannot have empty vm prefix")
	}

	if cmd.options.Stemecell == "" {
		return false, errors.New("cannot have empty stemcell")
	}

	return true, nil
}

func (cmd provisioningBaremetalCommand) Execute(args []string) (int, error) {
	_, err := cmd.printer.Printf("Executing %s command: args: %#v, options: %#v", cmd.Name(), args, cmd.options)
	if err != nil {
		return 1, err
	}

	_, err = cmd.ui.PrintlnInfo("WARNING: Be careful provisioning with the specific stemcell!")
	if err != nil {
		return 1, err
	}
	if cmd.isConfirmed() {
		provisioningBaremetalInfo := clients.ProvisioningBaremetalInfo{
			VmNamePrefix:     cmd.options.VMPrefix,
			Bm_stemcell:      cmd.options.Stemecell,
			Bm_netboot_image: cmd.options.NetbootImage,
		}

		provisioningBaremetalResponse, err := cmd.bmpClient.ProvisioningBaremetal(provisioningBaremetalInfo)
		if err != nil {
			return 1, err
		}

		if provisioningBaremetalResponse.Status != 200 {
			return provisioningBaremetalResponse.Status, nil
		}

		_, err = cmd.ui.PrintlnInfo("Provisioning Successful!")
		if err != nil {
			return 1, err
		}

		return 0, nil
	} else {
		return 1, errors.New("Provisioning Cancelled!")
	}
}

// Private Methods

func (cmd provisioningBaremetalCommand) isConfirmed() bool {
	var userInput string
	_, err := cmd.ui.PrintfInfo("Continue to provisioning? (type 'yes' to continue)")
	if err != nil {
		return false
	}
	_, err = cmd.ui.Scanln(&userInput)
	if err != nil {
		return false
	}

	if cmd.isYes(userInput) {
		return true
	}

	return false
}

func (cmd provisioningBaremetalCommand) isYes(userInput string) bool {
	yes := map[string]bool{
		"y":   true,
		"Y":   true,
		"yes": true,
		"Yes": true,
		"YES": true,
	}

	return yes[userInput]
}
