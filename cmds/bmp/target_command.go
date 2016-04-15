package bmp

import (
	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type targetCommand struct {
	args    []string
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewTargetCommand(options cmds.Options, bmpClient clients.BmpClient) targetCommand {
	consoleUi := common.NewConsoleUi()

	return targetCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd targetCommand) Name() string {
	return "target"
}

func (cmd targetCommand) Description() string {
	return "Set the URL of Bare Metal Provision Server"
}

func (cmd targetCommand) Usage() string {
	return "bmp target http://url.to.bmp.server"
}

func (cmd targetCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd targetCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: args: %#v, options: %#v", cmd.Name(), cmd.args, cmd.options)
	return true, nil
}

func (cmd targetCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)
	return 0, nil
}
