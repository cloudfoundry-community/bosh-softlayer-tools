package bmp

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	clients "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	cmds "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
	tablewriter "github.com/olekukonko/tablewriter"
)

type bmsCommand struct {
	args    []string
	options cmds.Options

	ui      common.UI
	printer common.Printer

	bmpClient clients.BmpClient
}

func NewBmsCommand(options cmds.Options, bmpClient clients.BmpClient) bmsCommand {
	consoleUi := common.NewConsoleUi(options.Verbose)

	return bmsCommand{
		options:   options,
		ui:        consoleUi,
		printer:   common.NewDefaultPrinter(consoleUi, options.Verbose),
		bmpClient: bmpClient,
	}
}

func (cmd bmsCommand) Name() string {
	return "bms"
}

func (cmd bmsCommand) Description() string {
	return "List all bare metals"
}

func (cmd bmsCommand) Usage() string {
	return "bmp bms --deployment[-d] <deployment file>"
}

func (cmd bmsCommand) Options() cmds.Options {
	return cmd.options
}

func (cmd bmsCommand) Validate() (bool, error) {
	cmd.printer.Printf("Validating %s command: args: %#v, options: %#v", cmd.Name(), cmd.args, cmd.options)
	if cmd.options.Deployment == "" {
		return false, errors.New("please specify the deployment file with -d")
	}

	_, err := os.Stat(cmd.options.Deployment)
	if os.IsNotExist(err) {
		return false, errors.New(fmt.Sprintf("deployment file '%s' doesn't exist", cmd.options.Deployment))
	}

	return true, nil
}

func (cmd bmsCommand) Execute(args []string) (int, error) {
	cmd.printer.Printf("Executing %s command: args: %#v, options: %#v", cmd.Name(), cmd.args, cmd.options)

	filename, _ := filepath.Abs(cmd.options.Deployment)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		errorMessage := fmt.Sprintf("bmp: could not read File '%s', error message '%s'", filename, err.Error())
		return 1, errors.New(errorMessage)
	}

	var deploymentInfo clients.Deployment
	err = yaml.Unmarshal(yamlFile, &deploymentInfo)
	if err != nil {
		errorMessage := fmt.Sprintf("bmp: failed to decode Yaml File '%s', error message '%s'", filename, err.Error())
		return 1, errors.New(errorMessage)
	}

	deploymentName := deploymentInfo.Name
	bmsResponse, err := cmd.bmpClient.Bms(deploymentName)
	if err != nil {
		return bmsResponse.Status, err
	}

	if bmsResponse.Status != 200 {
		return bmsResponse.Status, nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Hostname", "IPs", "Hardware_status", "Memory", "Cpu", "Provision_date"})

	content := make([][]string, len(bmsResponse.Data))
	table.SetHeader([]string{"Id", "Hostname", "IPs", "Hardware_status", "Memory", "Cpu", "Provision_date"})
	for i, serverInfo := range bmsResponse.Data {
		IPs := strings.Join([]string{serverInfo.Private_ip_address, serverInfo.Public_ip_address}, "/")
		content[i] = []string{
			strconv.Itoa(serverInfo.Id),
			serverInfo.Hostname,
			IPs,
			serverInfo.Hardware_status,
			strconv.Itoa(serverInfo.Memory),
			strconv.Itoa(serverInfo.Cpu),
			serverInfo.Provision_date}
	}

	for _, value := range content {
		table.Append(value)
	}
	table.Render()

	return 0, nil
}
