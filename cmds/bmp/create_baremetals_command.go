package bmp

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type createBaremetalsCommand struct {
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewCreateBaremetalsCommand(options cmds.Options, bmpClient clients.BmpClient) createBaremetalsCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return createBaremetalsCommand{
		options: options,

		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd createBaremetalsCommand) Name() string {
	return "create-baremetals"
}

func (cmd createBaremetalsCommand) Description() string {
	return `Create the missed baremetals: \"option --dryrun, only verify the orders\"`
}

func (cmd createBaremetalsCommand) Usage() string {
	return "bmp create-baremetals --deployment[-d] <deployment file> [--dryrun]"
}

func (cmd createBaremetalsCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd createBaremetalsCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: options: %#v", cmd.Name(), cmd.options)
	if cmd.options.Deployment == "" {
		return false, errors.New("please specify the deployment file with -d")
	}

	_, err := os.Stat(cmd.options.Deployment)
	if os.IsNotExist(err) {
		return false, errors.New(fmt.Sprintf("deployment file %s doesn't exist", cmd.options.Deployment))
	}

	return true, nil
}

func (cmd createBaremetalsCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s comamnd: args: %#v, options: %#v", cmd.Name(), args, cmd.options)

	filename, _ := filepath.Abs(cmd.options.Deployment)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not read File %s, error message %s", filename, err.Error())
		return 1, errors.New(errorMessage)
	}

	var deploymentInfo clients.Deployment
	err = yaml.Unmarshal(yamlFile, &deploymentInfo)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode Yaml File %s, error message %s", filename,
			err.Error())
		return 1, errors.New(errorMessage)
	}

	createBaremetalsInfo := clients.CreateBaremetalsInfo{}
	createBaremetalsInfo.Deployment = deploymentInfo.Name
	for _, resource := range deploymentInfo.ResourcePools {
		if resource.CloudProperty.Baremetal {
			createBaremetalsInfo.BaremetalSpecs = append(createBaremetalsInfo.BaremetalSpecs,
				resource.CloudProperty)
		}
	}

	createBaremetalsResponse, err := cmd.bmpClient.CreateBaremetals(createBaremetalsInfo, cmd.options.DryRun)
	if err != nil {
		return 1, err
	}

	if createBaremetalsResponse.Status != 200 {
		return createBaremetalsResponse.Status, nil
	}

	taskId := createBaremetalsResponse.Data.TaskId
	cmd.ui.Printf("Run command: bmp task %d to get the status of the task", taskId)

	return 0, nil
}
