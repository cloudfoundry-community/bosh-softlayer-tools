package stemcells

type LightStemcell interface {
	GetId() int
	GetUuid() string
	GetName() string
	GetDescription() string
	GetArchitecture() string
}
