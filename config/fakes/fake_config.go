package fakes

import (
	config "github.com/cloudfoundry-community/bosh-softlayer-tools/config"
)

type FakeConfig struct {
	Username, Password, TargetUrl string

	GetPathString string

	LoadConfigConfigInfo config.ConfigInfo
	LoadConfigErr        error

	SaveConfigConfigInfo config.ConfigInfo
	SaveConfigErr        error
}

func NewFakeConfig(username, password, targetUrl string) *FakeConfig {
	return &FakeConfig{
		Username:  username,
		Password:  password,
		TargetUrl: targetUrl,
	}
}

func (f *FakeConfig) GetPath() string {
	return f.GetPathString
}

func (f *FakeConfig) LoadConfig() (config.ConfigInfo, error) {
	return f.LoadConfigConfigInfo, f.LoadConfigErr
}

func (f *FakeConfig) SaveConfig(configInfo config.ConfigInfo) error {
	f.SaveConfigConfigInfo = configInfo

	return f.SaveConfigErr
}
