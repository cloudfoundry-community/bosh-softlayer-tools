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

// /bms
type BaremetalInfo struct {
	Id                 int    `json:"id"`
	Hostname           string `json:"hostname"`
	Private_ip_address string `json:"private_ip_address"`
	Public_ip_address  string `json:"public_ip_address"`
	Hardware_status    string `json:"hardware_status"`
	Memory             int    `json:"memory"`
	Cpu                int    `json:"cpu"`
	Provision_date     string `json:"provision_date"`
}

type BmsResponse struct {
	Status int             `json:"Status"`
	Data   []BaremetalInfo `json:"data"`
}

// /sl/packages

type Package struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type DataPackage struct {
	Packages []Package `json:"packages"`
}

type SlPackagesResponse struct {
	Status int         `json:"status"`
	Data   DataPackage `json:"data"`
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
	Status int                `json:"status"`
	Data   DataPackageOptions `json:"data"`
}

// /stemcells

type StemcellsResponse struct {
	Status   int      `json:"status"`
	Stemcell []string `json:"data"`
}

// /tasks?latest= (default 50)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	Status      string `json:"status"`
	EndTime     string `json:"end_time"`
}

type TasksResponse struct {
	Status int    `json:"status"`
	Data   []Task `json:"data"`
}

// /task/${task_id}/txt}" (default event)

type TaskOutputResponse struct {
	Status int      `json:"status"`
	Data   []string `json:"data`
}

// /baremetal/${serverId}/${status}

type UpdateStateResponse struct {
	Status int `json:"status"`
}

// /login/${username}/${password}

type LoginResponse struct {
	Status int `json:"status"`
}

// //baremetals (dry_run: optional)

type TaskInfo struct {
	TaskId int `json:"task_id"`
}

type CreateBaremetalsResponse struct {
	Status int      `json:"status"`
	Data   TaskInfo `json:"data"`
}

type ServerSpec struct {
	Package       string `yaml:"package"`
	Server        string `yaml:"server"`
	Ram           string `yaml:"ram"`
	Disk0         string `yaml:"disk0"`
	PortSpeed     string `yaml:"port_speed"`
	PublicVlanId  string `yaml:"public_vlan_id"`
	PrivateVlanId string `yaml:"private_vlan_id"`
	Hourly        bool   `yaml:"hourly"`
}

type CloudProperty struct {
	ImageId    string     `yaml:"image_id"`
	BoshIP     string     `yaml:"bosh_ip"`
	Datacenter string     `yaml:"datacenter"`
	NamePrefix string     `yaml:"name_prefix"`
	Baremetal  bool       `yaml:"baremetal"`
	ServerSpec ServerSpec `yaml:"server_spec"`
}

type CreateBaremetalsParameters struct {
	Parameters CreateBaremetalsInfo `json:"parameters"`
}

type CreateBaremetalsInfo struct {
	BaremetalSpecs []CloudProperty `json:"baremetal_specs"`
	Deployment     string          `json:"deployment"`
}

type Stemcell struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Resource struct {
	Name          string        `yaml:"name"`
	Network       string        `yaml:"network"`
	Size          uint          `yaml:"size"`
	Stemcell      Stemcell      `yaml:"stemcell"`
	CloudProperty CloudProperty `yaml:"cloud_properties"`
}

// deployment
type Deployment struct {
	Name          string     `yaml:"name"`
	ResourcePools []Resource `yaml:"resource_pools"`
}
