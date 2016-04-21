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
	Start_time  string `json:"start_time"`
	Status      string `json:"status"`
	End_time    string `json:"end_time"`
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

type UpdateStatusResponse struct {
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

type CreateBaremetalResponse struct {
	Status int      `json:"status"`
	Data   TaskInfo `json:"data"`
}

type ServerSpec struct {
	Package       string `json:"package"`
	Server        string `json:"server"`
	Ram           string `json:"ram"`
	Disk0         string `json:"disk0"`
	PortSpeed     string `json:"port_speed"`
	PublicVlanId  string `json:"public_vlan_id"`
	PrivateVlanId string `json:"private_vlan_id"`
	Hourly        bool   `json:"hourly"`
}

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

// deployment
type Deployment struct {
	Name string `json:"name"`
}
