package light_stemcell_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStemcells(t *testing.T) {
	BeforeEach(func() {
		os.Setenv("SL_USERNAME", "fake-username")
		username := os.Getenv("SL_USERNAME")
		Expect(username).To(Equal("fake-username"))

		os.Setenv("SL_API_KEY", "fake-api-key")
		apiKey := os.Getenv("SL_API_KEY")
		Expect(apiKey).To(Equal("fake-api-key"))
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Light Stemcell Suite")
}
