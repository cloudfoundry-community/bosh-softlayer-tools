package import_image_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestImportImage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Import Image Suite")
}
