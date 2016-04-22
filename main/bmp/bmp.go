package main

import (
	"errors"
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"

	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

var bmpOptions cmds.Options
var parser = flags.NewParser(&bmpOptions, flags.Default)

func main() {
	args, err := parser.ParseArgs(os.Args)
	if err != nil {
		fmt.Println("bmp: could not parse command args, err:", err)
		os.Exit(1)
	}

	command, err := createCommand(args, bmpOptions)
	if err != nil {
		fmt.Println("bmp: could not create command, err:", err)
		os.Exit(1)
	}

	validated, err := command.Validate()
	if err != nil {
		fmt.Println("bmp: could not validate command, err:", err)
		os.Exit(1)
	}

	if !validated {
		fmt.Println("bmp: invalid options for command: ", err.Error())
		os.Exit(1)
	}

	rc, err := command.Execute(args)
	if err != nil {
		fmt.Println("bmp: could not execute command, err:", err)
		os.Exit(rc)
	}

	os.Exit(rc)
}

func createCommand(args []string, options cmds.Options) (cmds.Command, error) {
	if len(args) < 2 {
		return nil, errors.New("No bmp command specified")
	}

	bmpClient, err := common.CreateBmpClient()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create BMP client: %s", err.Error()))
	}

	cmd := createCommands(options, bmpClient)[args[1]]
	if cmd == nil {
		return nil, errors.New(fmt.Sprintf("Invalid command: %s", args[1]))
	}

	return cmd, nil
}

func createCommands(options cmds.Options, bmpClient clients.BmpClient) map[string]cmds.Command {
	return map[string]cmds.Command{
		"bms":               bmp.NewBmsCommand(options, bmpClient),
		"create-baremetals": bmp.NewCreateBaremetalsCommand(options, bmpClient),
		"login":             bmp.NewLoginCommand(options, bmpClient),
		"sl":                bmp.NewSlCommand(options, bmpClient),
		"status":            bmp.NewStatusCommand(options, bmpClient),
		"stemcells":         bmp.NewStemcellsCommand(options, bmpClient),
		"target":            bmp.NewTargetCommand(options, bmpClient),
		"task":              bmp.NewTaskCommand(options, bmpClient),
		"tasks":             bmp.NewTasksCommand(options, bmpClient),
		"update-state":      bmp.NewUpdateStateCommand(options, bmpClient),
	}
}
