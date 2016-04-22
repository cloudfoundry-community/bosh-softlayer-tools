package bmp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/cloudfoundry-community/bosh-softlayer-tools/integration/test_helpers"
)

func TestBmp(t *testing.T) {
	BeforeSuite(test_helpers.BuildExecutables)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bmp Suite")
}
