package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var (
	DefaultLoader Loader = &localFileLoader{}
)

type LoaderConfig struct {
	LocalPath string `yaml:"local_path"`
	LocalFile string `yaml:"local_file"`
}

type Loader interface {
	Load(config LoaderConfig, dst *Config) (err error)
}

type localFileLoader struct{}

func (localFileLoader) Load(config LoaderConfig, dst *Config) (err error) {
	data, err := os.ReadFile(filepath.Join(config.LocalPath, config.LocalFile))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, dst)
	return
}
