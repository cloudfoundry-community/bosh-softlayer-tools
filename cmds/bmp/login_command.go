package bmp

import (
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type loginCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer
}

func NewLoginCommand(options cmds.Options) loginCommand {
	consoleUi := common.NewConsoleUi()

	return loginCommand{
		options: options,

		ui:      consoleUi,
		printer: common.NewDefaultPrinter(consoleUi, options.Verbose),
	}
}

func (cmd loginCommand) Name() string {
	return "login"
}

func (cmd loginCommand) Description() string {
	return "Login to the Bare Metal Provision Server"
}

func (cmd loginCommand) Usage() string {
	return "bmp login"
}

func (cmd loginCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd loginCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	return true, nil
}

func (cmd loginCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)
	return 0, nil
}
