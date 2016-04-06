package config

import (
	"encoding/json"
	"io/ioutil"
)

type ConfigInfo struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	TargetUrl string `json:"target_url"`
}

type Config interface {
	GetPath() string
	LoadConfig() (ConfigInfo, error)
	SaveConfig() error
}

const CONFIG_PATH = "~/.bmp_config"

type config struct {
	configInfo ConfigInfo
	path       string
}

func NewConfig(path string) *config {
	if path == "" {
		path = CONFIG_PATH
	}

	return &config{
		configInfo: ConfigInfo{},
		path:       path,
	}
}

func (c *config) GetPath() string {
	return c.path
}

func (c *config) LoadConfig() (ConfigInfo, error) {
	configFileContents, err := ioutil.ReadFile(c.path)
	if err != nil {
		return ConfigInfo{}, err
	}

	configInfo := ConfigInfo{}
	err = json.Unmarshal(configFileContents, &configInfo)
	if err != nil {
		return ConfigInfo{}, err
	}

	return configInfo, nil
}

func (c *config) SaveConfig() error {
	return nil
}
