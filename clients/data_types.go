package clients

// /info

type DataInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InfoResponse struct {
	Status int      `json:"status"`
	Data   DataInfo `json:"data"`
}

// /sl/packages

type DataPackage struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SlPackagesResponse struct {
	Data []DataPackage `json:"data"`
}

// /sl/${package_id}/options

type SlPackageOptionsResponse struct {
	//TODO
}

// /stemcells

type StemcellsResponse struct {
	//TODO
}
