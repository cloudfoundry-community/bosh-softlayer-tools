package bmp

import (
	"errors"

	"fmt"
	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type taskCommand struct {
	args    []string
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewTaskCommand(options cmds.Options, bmpClient clients.BmpClient) taskCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return taskCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd taskCommand) Name() string {
	return "task"
}

func (cmd taskCommand) Description() string {
	return `Show the output of the task: \"option --debug, Get the debug info of the task\"`
}

func (cmd taskCommand) Usage() string {
	return "bmp task <task-id>"
}

func (cmd taskCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd taskCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: args: %#v, options: %#v", cmd.Name(), cmd.args, cmd.options)

	if cmd.options.TaskID == 0 {
		return false, errors.New("cannot have empty task ID")
	}

	return true, nil
}

func (cmd taskCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)

	level := "event"
	if cmd.options.Debug == true {
		level = "deubg"
	}

	taskOutputResponse, err := cmd.bmpClient.TaskOutput(cmd.options.TaskID, level)
	if err != nil {
		return 1, err
	}

	if taskOutputResponse.Status != 200 {
		return taskOutputResponse.Status, nil
	}

	for _, value := range taskOutputResponse.Data {
		fmt.Println(value)
	}

	return 0, nil
}
