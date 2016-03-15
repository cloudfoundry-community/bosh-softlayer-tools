package cmds

import (
	"github.com/cloudfoundry-community/bosh-softlayer-tools/common"
)

type PrinterInterface interface {
	Println(a ...interface{}) (int, error)
	Printf(msg string, a ...interface{}) (int, error)
}

type CommandInterface interface {
	PrinterInterface
	Options() common.Options
	CheckOptions() error
	Run() error
}
