package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	slclient "github.com/maximilien/softlayer-go/client"
	"github.com/maximilien/softlayer-go/softlayer"

	"github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/import_image"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/cmds/light_stemcell"
	"github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

const VERSION = "v0.0.2"

var stemcellsOptions common.Options

type ImportImageCmdResponse struct {
	Id   int    `json:"id"`
	Uuid string `json:"uuid"`
}

func main() {
	flag.Parse()

	if stemcellsOptions.HelpFlag || stemcellsOptions.LongHelpFlag || len(os.Args) == 0 {
		usage()
		return
	}

	switch stemcellsOptions.CommandFlag {
	case "import-image":
		importImageCmd()
	case "light-stemcell":
		lightStemcellCmd()
	default:
		fmt.Println("SoftLayer BOSH Stemcells Utility")
	}
}

func importImageCmd() {
	if stemcellsOptions.HelpFlag {
		usage()
		return
	}

	client, err := createSoftLayerClient()
	if err != nil {
		fmt.Println("stemcells: Could not create the SoftLayer client, err:", err)
		os.Exit(1)
	}

	cmd, err := import_image.NewImportImageCmd(stemcellsOptions, client)
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
	common.TIMEOUT = 300 * time.Second
	common.POLLING_INTERVAL = 30 * time.Second

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
	if stemcellsOptions.HelpFlag {
		usage()
		return
	}

	client, err := createSoftLayerClient()
	if err != nil {
		fmt.Println("stemcells: Could not create the SoftLayer client, err:", err)
		os.Exit(1)
	}

	var cmd light_stemcell.LightStemcellCmd

	if stemcellsOptions.LightStemcellTypeFlag == "VDI" {
		cmd = light_stemcell.NewLightStemcellVDICmd(stemcellsOptions, client)
	} else {
		cmd = light_stemcell.NewLightStemcellVGBDGTCmd(stemcellsOptions, client)
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
	flag.StringVar(&stemcellsOptions.CommandFlag, "c", "", "the command, one of: import-image")

	flag.BoolVar(&stemcellsOptions.HelpFlag, "h", false, "prints the usage")
	flag.BoolVar(&stemcellsOptions.LongHelpFlag, "-help", false, "prints the usage")

	flag.BoolVar(&stemcellsOptions.PublicFlag, "public", false, "make the stemcell public")

	flag.StringVar(&stemcellsOptions.NameFlag, "name", "", "the name used by the specified command")
	flag.StringVar(&stemcellsOptions.NoteFlag, "note", "", "the note to be applied to the imported template")
	flag.StringVar(&stemcellsOptions.PublicNameFlag, "public-name", "", "the group name of public image to be imported")
	flag.StringVar(&stemcellsOptions.PublicNoteFlag, "public-note", "", "the note and summary of public image to be imported")
	flag.StringVar(&stemcellsOptions.OsRefCodeFlag, "os-ref-code", "UBUNTU_14_64", "the referenceCode of the operating system software")
	flag.StringVar(&stemcellsOptions.UriFlag, "uri", "", "the URI for an object storage object (.vhd/.iso file)")

	flag.StringVar(&stemcellsOptions.LightStemcellTypeFlag, "type", "VGBDGT", "two possible SoftLayer light stemcells: VGBDGT (default) or VDI")
	flag.StringVar(&stemcellsOptions.LightStemcellPathFlag, "path", ".", "the path for the location of the light stemcell file created")
	flag.StringVar(&stemcellsOptions.VersionFlag, "version", "", "the light stemcell version")
	flag.StringVar(&stemcellsOptions.StemcellInfoFilenameFlag, "stemcell-info-filename", "", "the path and filename to a JSON file containing the ID & UUID for a SoftLayer stemcell ")
	flag.StringVar(&stemcellsOptions.InfrastructureFlag, "infrastructure", "softlayer", "the light stemcell infrastructure, defaults to softlayer")
	flag.StringVar(&stemcellsOptions.HypervisorFlag, "hypervisor", "esxi", "the light stemcell version")
	flag.StringVar(&stemcellsOptions.OsNameFlag, "os-name", "ubuntu-trusty", "the name of the operating system")
}

func usage() {
	usageString := `
usage: sl-stemcells -c import-image [--name <template-name>] [--note <import note>]
       --os-ref-code <OsRefCode> --uri <swiftURI> --public [--public-name <public template-name>] [--public-note <public import note>]

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
  --public         the image will be made as public if this argument is specified
  --public-name    the group name of public image to be imported
  --public-note    the note and summary of public image to be imported
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
