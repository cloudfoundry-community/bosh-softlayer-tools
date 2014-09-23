package common_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	stemcells "github.com/maximilien/bosh-softlayer-stemcells/stemcells"
)

var _ = Describe("LightStemcellCmd", func() {
	var (
		lightStemcellMF stemcells.LightStemcellMF
	)

	Context("common.LoadStemcellMF", func() {
		It("loads and parses a new StemcellMF object from json file", func() {
			Fail("implement me!")
		})
	})
})