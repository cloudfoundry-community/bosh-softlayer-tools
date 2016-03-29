package bmp

import (
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type createBaremetalsCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer
}

func NewCreateBaremetalsCommand(options cmds.Options) createBaremetalsCommand {
	consoleUi := common.NewConsoleUi()

	return createBaremetalsCommand{
		options: options,

		ui:      consoleUi,
		printer: common.NewDefaultPrinter(consoleUi, options.Verbose),
	}
}

func (cmd createBaremetalsCommand) Name() string {
	return "create-baremetals"
}

func (cmd createBaremetalsCommand) Description() string {
	return `Create the missed baremetals: \"option --dryrun, only verify the orders\"`
}

func (cmd createBaremetalsCommand) Usage() string {
	return "bmp create-baremetals [--dryrun]"
}

func (cmd createBaremetalsCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd createBaremetalsCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	return true, nil
}

func (cmd createBaremetalsCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)
	return 0, nil
}
