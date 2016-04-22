package bmp

import (
	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type updateStateCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewUpdateStateCommand(options cmds.Options, bmpClient clients.BmpClient) updateStateCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return updateStateCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd updateStateCommand) Name() string {
	return "update-state"
}

func (cmd updateStateCommand) Description() string {
	return `Update the server state (\"bm.state.new\", \"bm.state.using\", \"bm.state.loading\", \"bm.state.failed\", \"bm.state.deleted\")`
}

func (cmd updateStateCommand) Usage() string {
	return "bmp update-state <bm.state.new>"
}

func (cmd updateStateCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd updateStateCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	return true, nil
}

func (cmd updateStateCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s command: args: %#v, options: %#v", cmd.Name(), args, cmd.options)
	return 0, nil
}
