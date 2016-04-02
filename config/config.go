package config

type ConfigInfo struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	TargetUrl string `json:"target_url"`
}

type Config interface {
	GetPath() string
	LoadConfig() (Config, error)
	SaveConfig() error
}

const CONFIG_PATH = "~/.bmp_config"

type config struct {
	configInfo ConfigInfo
	path       string
}

func NewConfig(path string) *config {
	return &config{
		configInfo: ConfigInfo{},
		path:       path,
	}
}

func (c *config) GetPath() string {
	return ""
}

func (c *config) LoadConfig() (ConfigInfo, error) {
	return ConfigInfo{}, nil
}

func (c *config) SaveConfig() error {
	return nil
}
