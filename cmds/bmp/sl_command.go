package bmp

import (
	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type slCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewSlCommand(options cmds.Options, bmpClient clients.BmpClient) slCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return slCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd slCommand) Name() string {
	return "sl"
}

func (cmd slCommand) Description() string {
	return "List all Softlayer packages or package options"
}

func (cmd slCommand) Usage() string {
	return "bmp sl --packages | --package-options"
}

func (cmd slCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd slCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	return true, nil
}

func (cmd slCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)
	return 0, nil
}
