package clients

type BmpClient interface {
	Info() (InfoResponse, error)
	SlPackages() (SlPackagesResponse, error)

	// Stemcells() (StemcellsResponse, error)
	// SlPackageOptions(packageId string) (SlPackageOptionsResponse, error)
}
