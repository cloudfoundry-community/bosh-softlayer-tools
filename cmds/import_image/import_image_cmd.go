package import_image

import (
	"fmt"

	"github.com/maximilien/bosh-softlayer-stemcells/common"
)

type importImageCmd struct {
	options common.Options

	Name      string
	Note      string
	OsRefCode string
	Uri       string
}

func NewImportImageCmd(options common.Options) (importImageCmd, error) {
	return importImageCmd{
		options: options,

		Name:      options.NameFlag,
		Note:      options.NoteFlag,
		OsRefCode: options.OsRefCodeFlag,
		Uri:       options.UriFlag,
	}, nil
}

func (cmd importImageCmd) Println(a ...interface{}) (int, error) {
	fmt.Println(a)

	return 0, nil
}

func (cmd importImageCmd) Printf(msg string, a ...interface{}) (int, error) {
	fmt.Printf(msg, a)

	return 0, nil
}

func (cmd importImageCmd) Options() common.Options {
	return cmd.options
}

func (cmd importImageCmd) Run() error {
	cmd.Println("Implment me!: Run")

	return nil
}
