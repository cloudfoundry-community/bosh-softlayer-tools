package bmp

import (
	"errors"
	"os"
	"strconv"

	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
	"github.com/olekukonko/tablewriter"
)

type slCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewSlCommand(options cmds.Options, bmpClient clients.BmpClient) slCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return slCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd slCommand) Name() string {
	return "sl"
}

func (cmd slCommand) Description() string {
	return "List all Softlayer packages or package options"
}

func (cmd slCommand) Usage() string {
	return "bmp sl --packages | --package-options"
}

func (cmd slCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd slCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)

	if !cmd.options.Packages && cmd.options.PackageOptions == "" {
		return false, errors.New("Please specify --packages or --package-options")
	}

	return true, nil
}

func (cmd slCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)

	validate, err := cmd.Validate()
	if validate == false && err != nil {
		return 1, errors.New("bmp sl validation error")
	}

	if cmd.options.Packages {
		slPackagesResponse, err := cmd.bmpClient.SlPackages()
		if err != nil {
			return 1, err
		}

		if slPackagesResponse.Status != 200 {
			return slPackagesResponse.Status, nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Package ID", "Name"})

		length := len(slPackagesResponse.Data.Packages)
		content := make([][]string, length)
		for i, slPackage := range slPackagesResponse.Data.Packages {
			content[i] = []string{
				strconv.Itoa(slPackage.Id),
				slPackage.Name}
		}

		for _, value := range content {
			table.Append(value)
		}

		cmd.ui.PrintTable(table)
		cmd.ui.Println("Packages total:", length)

		return 0, nil
	} else {
		//TODO: Add sl --package-options

		return 0, nil
	}
}
