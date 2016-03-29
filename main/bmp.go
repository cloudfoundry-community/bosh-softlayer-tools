package main

import (
	"errors"
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"

	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	bmp "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/bmp"
)

var options cmds.Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	args, err := parser.ParseArgs(os.Args)
	if err != nil {
		fmt.Println("bmp: could not parse command args, err:", err)
		os.Exit(1)
	}

	command, err := createCommand(args, options)
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

	cmd := createCommands(options)[args[1]]
	if cmd == nil {
		return nil, errors.New(fmt.Sprintf("Invalid command: %s", args[1]))
	}

	return cmd, nil
}

func createCommands(options cmds.Options) map[string]cmds.Command {
	return map[string]cmds.Command{
		"bms":               bmp.NewBmsCommand(options),
		"create-baremetals": bmp.NewCreateBaremetalsCommand(options),
		"login":             bmp.NewLoginCommand(options),
		"sl":                bmp.NewSlCommand(options),
		"status":            bmp.NewStatusCommand(options),
		"stemcells":         bmp.NewStemcellsCommand(options),
		"target":            bmp.NewTargetCommand(options),
		"task":              bmp.NewTaskCommand(options),
		"tasks":             bmp.NewTasksCommand(options),
		"update-state":      bmp.NewUpdateStateCommand(options),
	}
}
