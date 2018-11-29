package bmp

import (
	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type statusCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewStatusCommand(options cmds.Options, bmpClient clients.BmpClient) statusCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return statusCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd statusCommand) Name() string {
	return "status"
}

func (cmd statusCommand) Description() string {
	return "show bmp status"
}

func (cmd statusCommand) Usage() string {
	return "bmp status"
}

func (cmd statusCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd statusCommand) Validate() (bool, error) {
	_, err := cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (cmd statusCommand) Execute(args []string) (int, error) {
	_, err := cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)
	if err != nil {
		return 1, err
	}

	info, err := cmd.bmpClient.Info()
	if err != nil {
		return 1, err
	}

	_, err = cmd.ui.PrintlnInfo("BMP server info")
	if err != nil {
		return 1, err
	}
	_, err = cmd.ui.PrintfInfo(" name:    %s\n", info.Data.Name)
	if err != nil {
		return 1, err
	}
	_, err = cmd.ui.PrintfInfo(" version: %s\n", info.Data.Version)
	if err != nil {
		return 1, err
	}

	return 0, nil
}
