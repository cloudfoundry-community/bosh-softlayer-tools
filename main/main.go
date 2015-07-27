package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"runtime/debug"

	slclient "github.com/maximilien/softlayer-go/client"
	softlayer "github.com/maximilien/softlayer-go/softlayer"

	import_image "github.com/maximilien/bosh-softlayer-stemcells/cmds/import_image"
	common "github.com/maximilien/bosh-softlayer-stemcells/common"
)

const VERSION = "v0.0.1"

var options common.Options

func main() {
	defer handlePanic()

	if options.HelpFlag || options.LongHelpFlag || len(flag.Args()) == 0 {
		usage()
		return
	}

	switch flag.Args()[0] {
	case "import-image":
		importImageCmd()
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

	startTime := time.Now()

	err = cmd.Run()
	if err != nil {
		cmd.Println("stemcells: Could not import image, err:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	cmd.Println("Total time:", duration)
}

func init() {
	flag.BoolVar(&options.HelpFlag, "h", false, "prints the usage")
	flag.BoolVar(&options.LongHelpFlag, "-help", false, "prints the usage")

	flag.StringVar(&options.NameFlag, "name", "stemcell-image", "the group name to be applied to the imported template")
	flag.StringVar(&options.NoteFlag, "note", "", "the note to be applied to the imported template")
	flag.StringVar(&options.OsRefCodeFlag, "os-ref-code", "UBUNTU_14_64", "the referenceCode of the operating system software")
	flag.StringVar(&options.UriFlag, "uri", "", "the URI for an object storage object (.vhd/.iso file)")

	flag.Parse()
}

func usage() {
	usageString := `
  usage: bosh-softlayer-stemcells import-image [-v] [--dry-run]

  -h | --help     prints the usage

  IMPORT-IMAGE:

  import-image 	  the import image command
  --name 			the group name to be applied to the imported template
  --note 			the note to be applied to the imported template
  --os-ref-code	  the referenceCode of the operating system software 
				  description for the imported VHD 
				  available options: CENTOS_6_32, CENTOS_6_64, CENTOS_7_64, 
					REDHAT_6_32, REDHAT_6_64, REDHAT_7_64, UBUNTU_10_32, 
					UBUNTU_10_64, UBUNTU_12_32, UBUNTU_12_64, UBUNTU_14_32, 
					UBUNTU_14_64, WIN_2003-STD-SP2-5_32, WIN_2003-STD-SP2-5_64, 
					WIN_2012-STD_64
  --uri 		  the URI for an object storage object (.vhd/.iso file)
				  swift://<ObjectStorageAccountName>@<clusterName>/<containerName>/<fileName.(vhd|iso)>
	`

	fmt.Println(fmt.Sprintf("%s\nVersion %s", usageString, VERSION))
}

func handlePanic() {
	err := recover()
	if err != nil {
		switch err := err.(type) {
		case error:
			displayCrashDialog(err.Error())
		case string:
			displayCrashDialog(err)
		default:
			displayCrashDialog("An unexpected type of error")
		}
	}

	if err != nil {
		os.Exit(1)
	}
}

func displayCrashDialog(errorMessage string) {
	formattedString := `
Something completely unexpected happened. This is a bug in %s.
Please file this bug : https://github.com/maximilien/bosh-softlayer-stemcells/issues
Tell us that you ran this command:

	%s

this error occurred:

	%s

and this stack trace:

%s
	`

	stackTrace := "\t" + strings.Replace(string(debug.Stack()), "\n", "\n\t", -1)
	println(fmt.Sprintf(formattedString, "bosh-softlayer-stemcells", strings.Join(os.Args, " "), errorMessage, stackTrace))
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
