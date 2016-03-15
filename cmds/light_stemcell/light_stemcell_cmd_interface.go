package light_stemcell

import (
	cmds "github.com/cloudfoundry-incubator/bosh-softlayer-tools/cmds"
)

type SoftLayerStemcellInfo struct {
	Id             int    `json:"id"`
	Uuid           string `json:"uuid"`
	DatacenterName string `json:"datacenter-name"`
}

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

type LightStemcellMF struct {
	Name            string          `json:"name" yaml:"name"`
	Version         string          `json:"version" yaml:"version"`
	BoshProtocol    int             `json:"bosh_protocol" yaml:"bosh_protocol"`
	Sha1            string          `json:"sha1" yaml:"sha1"`
	CloudProperties CloudProperties `json:"cloud_properties" yaml:"cloud_properties"`
}

type CloudProperties struct {
	Infrastructure string `json:"infrastructure" yaml:"infrastructure"`
	Architecture   string `json:"architecture" yaml:"architecture"`
	RootDeviceName string `json:"root_device_name" yaml:"root_device_name"`

	//SoftLayer-specific properties
	VirtualDiskImageId   int    `json:"virtual-disk-image-id" yaml:"virtual-disk-image-id"`
	VirtualDiskImageUuid string `json:"virtual-disk-image-uuid" yaml:"virtual-disk-image-uuid"`
	DatacenterName       string `json:"datacenter-name" yaml:"datacenter-name"`
}

type LightStemcellCmd interface {
	cmds.CommandInterface

	GetStemcellPath() string
	GetLightStemcellInfo() LightStemcellInfo
	Create(imageId int) (string, error)
}
