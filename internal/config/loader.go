package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	DefaultLoader Loader = &localFileLoader{}
)

type LoaderConfig struct {
	LocalPath string `json:"local_path"`
	LocalFile string `json:"local_file"`
}

type Loader interface {
	Load(config LoaderConfig, dst interface{}) (err error)
}

type localFileLoader struct{}

func (localFileLoader) Load(config LoaderConfig, dst interface{}) (err error) {
	config.LocalPath = strings.Trim(config.LocalPath, "/")
	if config.LocalPath == "" {
		config.LocalPath = baseDir()
	}
	if config.LocalFile == "" {
		config.LocalFile = "config.json"
	}
	data, err := os.ReadFile(filepath.Join(config.LocalPath, config.LocalFile))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &dst)
	return
}

func baseDir() string {
	_, f, _, _ := runtime.Caller(1)
	return strings.TrimRight(filepath.Dir(f), "/")
}
