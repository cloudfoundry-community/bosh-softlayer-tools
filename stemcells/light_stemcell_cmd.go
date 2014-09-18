package stemcells

import (
	softlayer "github.com/maximilien/softlayer-go/softlayer"
)

type lightStemcell struct {
	Id           int
	Uuid         string
	Name         string
	Description  string
	Architecture string

	client softlayer.Client
}

func NewLightStemcellCmd(client softlayer.Client, vdImageId int) LightStemcell {
	return &lightStemcell{
		Id: vdImageId,

		client: client,
	}
}

func (ls *lightStemcell) GetId() int {
	return ls.Id
}

func (ls *lightStemcell) GetUuid() string {
	return ls.Uuid
}

func (ls *lightStemcell) GetName() string {
	return ls.Name
}

func (ls *lightStemcell) GetDescription() string {
	return ls.Description
}

func (ls *lightStemcell) GetArchitecture() string {
	return ls.Architecture
}
