package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
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
	config.LocalPath = strings.Trim(config.LocalPath, "/")
	if config.LocalPath == "" {
		config.LocalPath = baseDir()
	}
	if config.LocalFile == "" {
		config.LocalFile = "config.yaml"
	}
	data, err := os.ReadFile(filepath.Join(config.LocalPath, config.LocalFile))
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, dst)
	return
}

func baseDir() string {
	_, f, _, _ := runtime.Caller(1)
	return strings.TrimRight(filepath.Dir(filepath.Dir(filepath.Dir(f))), "/")
}
