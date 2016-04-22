package bmp

import (
	"errors"
	"net/url"

	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
	config "github.com/cloudfoundry-community/bosh-softlayer-tools/config"
)

type targetCommand struct {
	args    []string
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewTargetCommand(options cmds.Options, bmpClient clients.BmpClient) targetCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

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

	if cmd.options.Target == "" {
		return false, nil
	}

	_, err := url.ParseRequestURI(cmd.options.Target)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (cmd targetCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)

	validate, err := cmd.Validate()
	if validate == false && err == nil {
		return 1, errors.New("bmp target validation error")
	} else if validate == false && err != nil {
		return 1, err
	}

	configInfo, err := common.CreateConfig(cmd.bmpClient.ConfigPath())
	if err != nil {
		return 1, err
	}

	configInfo.TargetUrl = cmd.options.Target

	c := config.NewConfig(cmd.bmpClient.ConfigPath())
	err = c.SaveConfig(configInfo)
	if err != nil {
		return 1, err
	}

	return 0, nil
}
