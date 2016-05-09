package bmp

import (
	"os"
	"strconv"

	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
	tablewriter "github.com/olekukonko/tablewriter"
)

type tasksCommand struct {
	args    []string
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewTasksCommand(options cmds.Options, bmpClient clients.BmpClient) tasksCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return tasksCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd tasksCommand) Name() string {
	return "tasks"
}

func (cmd tasksCommand) Description() string {
	return `List tasks: \"option --latest number\", \"The number of latest tasks, default is 50\"`
}

func (cmd tasksCommand) Usage() string {
	return "bmp tasks"
}

func (cmd tasksCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd tasksCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: args: %#v, options: %#v", cmd.Name(), cmd.args, cmd.options)
	return true, nil
}

func (cmd tasksCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)

	tasksResponse, err := cmd.bmpClient.Tasks(cmd.options.Latest)
	if err != nil {
		return 1, err
	}

	if tasksResponse.Status != 200 {
		return tasksResponse.Status, nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Task ID", "Status", "Description", "Start", "End"})

	length := len(tasksResponse.Data)
	content := make([][]string, length)
	for i, task := range tasksResponse.Data {
		content[i] = []string{
			strconv.Itoa(task.Id),
			task.Description,
			task.Status,
			task.StartTime,
			task.EndTime}
	}

	for _, value := range content {
		table.Append(value)
	}

	cmd.ui.Println(table)
	cmd.ui.Println("Tasks total:", length)
	return 0, nil
}
