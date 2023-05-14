package config

import (
	"os"
	"sync"
)

var (
	once   sync.Once
	config Config
)

func GetConfig() Config {
	once.Do(func() {
		loaderConfig := LoaderConfig{
			LocalPath: os.Getenv("LOCAL_PATH"),
			LocalFile: os.Getenv("LOCAL_FILE"),
		}
		if err := DefaultLoader.Load(loaderConfig, &config); err != nil {
			panic(err)
		}
	})
	return config
}

// GetConfigApplication
// @autowire(set=config)
func GetConfigApplication() Application {
	return GetConfig().Application
}

// GetConfigBackendList
// @autowire(set=config)
func GetConfigBackendList() BackendList {
	list := GetConfig().BackendList
	list.Default()
	return list
}
