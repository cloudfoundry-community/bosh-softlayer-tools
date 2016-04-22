package test_helpers

import (
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

var BmpExec string
var StemcellsExec string

func BuildExecutables() {
	var err error
	BmpExec, err = gexec.Build("./../../main/bmp")
	Expect(err).NotTo(HaveOccurred())

	StemcellsExec, err = gexec.Build("./../../main/stemcells")
	Expect(err).NotTo(HaveOccurred())
}
