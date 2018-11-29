package bmp

import (
	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type stemcellsCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewStemcellsCommand(options cmds.Options, bmpClient clients.BmpClient) stemcellsCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return stemcellsCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd stemcellsCommand) Name() string {
	return "stemcells"
}

func (cmd stemcellsCommand) Description() string {
	return "List all stemcells"
}

func (cmd stemcellsCommand) Usage() string {
	return "bmp stemcells"
}

func (cmd stemcellsCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd stemcellsCommand) Validate() (bool, error) {
	_, err := cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (cmd stemcellsCommand) Execute(args []string) (int, error) {
	_, err := cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)
	if err != nil {
		return 1, err
	}

	stemcellsResponse, err := cmd.bmpClient.Stemcells()
	if err != nil {
		return 1, err
	}

	if stemcellsResponse.Status != 200 {
		return stemcellsResponse.Status, nil
	}

	table := cmd.ui.NewTableWriter()
	table.SetHeader([]string{"Stemcell"})
	length := len(stemcellsResponse.Stemcell)
	content := make([][]string, length)
	for i, stemcell := range stemcellsResponse.Stemcell {
		content[i] = []string{stemcell}
	}

	for _, value := range content {
		table.Append(value)
	}

	_, err = cmd.ui.PrintTable(table)
	if err != nil {
		return 1, err
	}
	_, err = cmd.ui.PrintlnInfo("")
	if err != nil {
		return 1, err
	}
	_, err = cmd.ui.PrintfInfo("Stemcells total: %d\n", length)
	if err != nil {
		return 1, err
	}

	return 0, nil
}
