package bmp

import (
	"errors"

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
	return `Show the output of the task: \"option --debug, Get the debug info of the task; --json, show info with JSON format\"`
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
		level = "debug"
	}

	if cmd.options.JSON {
		taskJsonResponse, err := cmd.bmpClient.TaskJsonOutput(cmd.options.TaskID, level)
		if err != nil {
			return 1, err
		}

		if taskJsonResponse.Status != 200 {
			return taskJsonResponse.Status, nil
		}

		cmd.ui.PrintfInfo("Task output for ID %d with %s level\n", cmd.options.TaskID, level)
		for _, value := range taskJsonResponse.Data {
			cmd.ui.PrintlnInfo(value)
		}

		return 0, nil

	} else {
		taskTxtResponse, err := cmd.bmpClient.TaskOutput(cmd.options.TaskID, level)
		if err != nil {
			return 1, err
		}

		if taskTxtResponse.Status != 200 {
			return taskTxtResponse.Status, nil
		}

		cmd.ui.PrintfInfo("Task output for ID %d with %s level\n", cmd.options.TaskID, level)
		for _, value := range taskTxtResponse.Data {
			cmd.ui.PrintlnInfo(value)
		}

		return 0, nil
	}
}
