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

// GetBackendList
// @autowire(set=config)
func GetBackendList() BackendSlice {
	return GetConfig().BackendList
}
