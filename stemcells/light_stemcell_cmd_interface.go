package stemcells

type LightStemcellInfo struct {
	//Defaulted
	Infrastructure string `json:"infrastructure"`
	Architecture   string `json:"architecture"`
	RootDeviceName string `json:"root-device-name"`

	//Required
	Version    string `json:"version"`
	Hypervisor string `json:"hypervisor"`
	OsName     string `json:"os-name"`
}

type LightStemcellCmd interface {
	GenerateStemcellName(info LightStemcellInfo) string
	GetStemcellsPath() string
	GetLightStemcellInfo() LightStemcellInfo
	Create(vdImageId int) (string, error)
}
