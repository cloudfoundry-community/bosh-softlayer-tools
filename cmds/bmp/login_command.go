package bmp

import (
	"errors"

	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
	config "github.com/cloudfoundry-community/bosh-softlayer-tools/config"
)

type loginCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewLoginCommand(options cmds.Options, bmpClient clients.BmpClient) loginCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return loginCommand{
		options: options,

		ui:      consoleUi,
		printer: common.NewDefaultPrinter(consoleUi, options.Verbose),

		bmpClient: bmpClient,
	}
}

func (cmd loginCommand) Name() string {
	return "login"
}

func (cmd loginCommand) Description() string {
	return "Login to the Bare Metal Provision Server"
}

func (cmd loginCommand) Usage() string {
	return "bmp login --username[-u] <username> --password[-p] <password>"
}

func (cmd loginCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd loginCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	if cmd.options.Username == "" {
		return false, errors.New("cannot have empty username")
	}

	if cmd.options.Password == "" {
		return false, errors.New("cannot have empty password")
	}

	return true, nil
}

func (cmd loginCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)

	loginResponse, err := cmd.bmpClient.Login(cmd.options.Username, cmd.options.Password)
	if err != nil {
		return loginResponse.Status, err
	}

	if loginResponse.Status != 200 {
		return loginResponse.Status, nil
	}

	c := config.NewConfig(cmd.bmpClient.ConfigPath())
	configInfo, err := c.LoadConfig()
	if err != nil {
		return 1, err
	}

	configInfo.Username = cmd.options.Username
	configInfo.Password = cmd.options.Password

	err = c.SaveConfig(configInfo)
	if err != nil {
		return 1, err
	}

	return 0, nil
}
