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
type Option struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type Category struct {
	Code     string   `json:"code"`
	Name     string   `json:"name"`
	Options  []Option `json:"options"`
	Required bool     `json:"required"`
}

type DataPackageOptions struct {
	Category   []Category `json:"categories"`
	Datacenter []string   `json:"datacenters"`
}
type SlPackageOptionsResponse struct {
	Data DataPackageOptions `json:"data"`
}

// /stemcells

type StemcellsResponse struct {
	Stemcell []string `json:"data"`
}
