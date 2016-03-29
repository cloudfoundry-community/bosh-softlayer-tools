package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	slclient "github.com/maximilien/softlayer-go/client"
	softlayer "github.com/maximilien/softlayer-go/softlayer"

	import_image "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/import_image"
	light_stemcell "github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/light_stemcell"
	common "github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

const VERSION = "v0.0.1"

var options common.Options

type ImportImageCmdResponse struct {
	Id   int    `json:"id"`
	Uuid string `json:"uuid"`
}

func main() {
	flag.Parse()

	if options.HelpFlag || options.LongHelpFlag || len(flag.Args()) == 0 {
		usage()
		return
	}

	switch flag.Args()[0] {
	case "import-image":
		importImageCmd()
	case "light-stemcell":
		lightStemcellCmd()
	default:
		fmt.Println("SoftLayer BOSH Stemcells Utility")
	}
}

func importImageCmd() {
	if options.HelpFlag {
		usage()
		return
	}

	client, err := createSoftLayerClient()
	if err != nil {
		fmt.Println("stemcells: Could not create the SoftLayer client, err:", err)
		os.Exit(1)
	}

	cmd, err := import_image.NewImportImageCmd(options, client)
	if err != nil {
		cmd.Println("stemcells: Could not create image command, err:", err)
		os.Exit(1)
	}

	err = cmd.CheckOptions()
	if err != nil {
		cmd.Println("stemcells: Could not create image command, err:", err)
		usage()
		os.Exit(1)
	}

	startTime := time.Now()

	err = cmd.Run()
	if err != nil {
		cmd.Println("stemcells: Could not import image, err:", err)
		os.Exit(1)
	}

	cmdOutput := &ImportImageCmdResponse{
		Id:   cmd.Id,
		Uuid: cmd.Uuid}

	cmdOutputJson, err := json.Marshal(cmdOutput)
	if err != nil {
		cmd.Println("Cannot marshal: improperly formatted json:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	cmd.Println("Total time: ", duration)

	fmt.Println(string(cmdOutputJson))
}

func lightStemcellCmd() {
	if options.HelpFlag {
		usage()
		return
	}

	client, err := createSoftLayerClient()
	if err != nil {
		fmt.Println("stemcells: Could not create the SoftLayer client, err:", err)
		os.Exit(1)
	}

	var cmd light_stemcell.LightStemcellCmd

	if options.LightStemcellTypeFlag == "VDI" {
		cmd = light_stemcell.NewLightStemcellVDICmd(options, client)
	} else {
		cmd = light_stemcell.NewLightStemcellVGBDGTCmd(options, client)
	}

	err = cmd.CheckOptions()
	if err != nil {
		cmd.Println("stemcells: Could not light stemcell command, err:", err)
		os.Exit(1)
	}

	startTime := time.Now()

	err = cmd.Run()
	if err != nil {
		cmd.Println("stemcells: Could not light stemcell, err:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	cmd.Println("Total time: ", duration)

	fmt.Println(cmd.GetStemcellPath())
}

func init() {
	flag.StringVar(&options.CommandFlag, "c", "", "the command, one of: import-image")

	flag.BoolVar(&options.HelpFlag, "h", false, "prints the usage")
	flag.BoolVar(&options.LongHelpFlag, "-help", false, "prints the usage")

	flag.StringVar(&options.NameFlag, "name", "", "the name used by the specified command")
	flag.StringVar(&options.NoteFlag, "note", "", "the note to be applied to the imported template")
	flag.StringVar(&options.OsRefCodeFlag, "os-ref-code", "UBUNTU_14_64", "the referenceCode of the operating system software")
	flag.StringVar(&options.UriFlag, "uri", "", "the URI for an object storage object (.vhd/.iso file)")

	flag.StringVar(&options.LightStemcellTypeFlag, "type", "VGBDGT", "two possible SoftLayer light stemcells: VGBDGT (default) or VDI")
	flag.StringVar(&options.LightStemcellPathFlag, "path", ".", "the path for the location of the light stemcell file created")
	flag.StringVar(&options.VersionFlag, "version", "", "the light stemcell version")
	flag.StringVar(&options.StemcellInfoFilenameFlag, "stemcell-info-filename", "", "the path and filename to a JSON file containing the ID & UUID for a SoftLayer stemcell ")
	flag.StringVar(&options.InfrastructureFlag, "infrastructure", "softlayer", "the light stemcell infrastructure, defaults to softlayer")
	flag.StringVar(&options.HypervisorFlag, "hypervisor", "esxi", "the light stemcell version")
	flag.StringVar(&options.OsNameFlag, "os-name", "ubuntu-trusty", "the name of the operating system")
}

func usage() {
	usageString := `
usage: bosh-softlayer-stemcells -c import-image [--name <template-name>] [--note <import note>] 
       --os-ref-code <OsRefCode> --uri <swiftURI>

  -h | --help   prints the usage

  IMPORT-IMAGE:

  -c import-image  the import image command
  --name           the group name to be applied to the imported template
  --note           the note to be applied to the imported template
  --os-ref-code    the referenceCode of the operating system software 
                   description for the imported VHD 
                   available options: CENTOS_6_32, CENTOS_6_64, CENTOS_7_64, 
                     REDHAT_6_32, REDHAT_6_64, REDHAT_7_64, UBUNTU_10_32, 
                     UBUNTU_10_64, UBUNTU_12_32, UBUNTU_12_64, UBUNTU_14_32, 
                     UBUNTU_14_64, WIN_2003-STD-SP2-5_32, WIN_2003-STD-SP2-5_64, 
                     WIN_2012-STD_64
  --uri            the URI for an object storage object (.vhd/.iso file)
                   swift://<ObjectStorageAccountName>@<clusterName>/<containerName>/<fileName.(vhd|iso)>
    `

	fmt.Println(fmt.Sprintf("%s\nVersion %s", usageString, VERSION))
}

func createSoftLayerClient() (softlayer.Client, error) {
	username := os.Getenv("SL_USERNAME")
	if username == "" {
		return nil, errors.New("stemcells: cannot create SoftLayer client: SL_USERNAME environment variable must be set")
	}

	apiKey := os.Getenv("SL_API_KEY")
	if apiKey == "" {
		return nil, errors.New("stemcells: cannot create SoftLayer client: SL_API_KEY environment variable must be set")
	}

	return slclient.NewSoftLayerClient(username, apiKey), nil
}
