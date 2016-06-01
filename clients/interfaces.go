package clients

type BmpClient interface {
	ConfigPath() string

	Info() (InfoResponse, error)
	Bms(deploymentName string) (BmsResponse, error)
	SlPackages() (SlPackagesResponse, error)
	Stemcells() (StemcellsResponse, error)
	SlPackageOptions(packageId string) (SlPackageOptionsResponse, error)
	TaskOutput(taskId uint, level string) (TaskOutputResponse, error)
	Tasks(latest uint) (TasksResponse, error)
	UpdateState(serverId string, status string) (UpdateStateResponse, error)
	Login(username string, password string) (LoginResponse, error)
	CreateBaremetals(createBaremetalsInfo CreateBaremetalsInfo, dryRun bool) (CreateBaremetalsResponse, error)
}
