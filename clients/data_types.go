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

// /tasks?latest= (default 50)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Start_time  string `json:"start_time"`
	Status      string `json:"status"`
	End_time    string `json:"end_time"`
}

type TasksResponse struct {
	Data []Task `json:"data"`
}

// /task/${task_id}/txt}" (default event)

type TaskOutputResponse struct {
	Data []string `json:"data`
}

// /baremetal/${serverId}/${status}
type UpdateStatusResponse struct {
	Status string `json:"status"`
}

// /login/${username}/${password}
type LoginResponse struct {
	Status string `json:"status"`
}

// //baremetals (dry_run: optional)
type TaskInfo struct {
	TaskId int `json:"task_id"`
}

type CreateBaremetalResponse struct {
	Data TaskInfo `json:"data"`
}

type ServerSpec struct {
	Package       string `json:"package"`
	Server        string `json:"server"`
	Ram           string `json:"ram"`
	Disk0         string `json:"disk0"`
	PortSpeed     string `json:"port_speed"`
	PublicVlanId  string `json:"public_vlan_id"`
	PrivateVlanId string `json:"private_vlan_id"`
	Hourly        bool `json:"hourly"`
}

//type Deployment struct {
//	Name          string         `yaml:"name"`
//	ResourcePools []ResourcePool `yaml:"resource_pools"`
//}

//type ResourcePool struct {
//	CloudProperties CloudProperty `yaml:"cloud_properties"`
//	Size            int           `yaml:"size"`
//}

type CloudProperty struct {
	ImageId    string     `json:"image_id"`
	BoshIP     string     `json:"bosh_ip"`
	Datacenter string     `json:"datacenter"`
	NamePrefix string     `json:"name_prefix"`
	Baremetal  bool       `json:"baremetal"`
	ServerSpec ServerSpec `json:"server_spec"`
}

type CreateBaremetalParameters struct {
	Parameters CreateBaremetalInfo `json:"parameters"`
}

type CreateBaremetalInfo struct {
	BaremetalSpecs []CloudProperty `json:"baremetal_specs"`
	Deployment     string          `json:"deployment"`
}
