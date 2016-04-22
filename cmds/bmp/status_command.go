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
	cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	return true, nil
}

func (cmd statusCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)

	info, err := cmd.bmpClient.Info()
	if err != nil {
		return 1, err
	}

	cmd.ui.Println("BMP server info")
	cmd.ui.Printf(" status: %d", info.Status)
	cmd.ui.Printf(" name: %s", info.Data.Name)
	cmd.ui.Printf(" version: %s", info.Data.Version)
	cmd.ui.Println()

	return 0, nil
}
