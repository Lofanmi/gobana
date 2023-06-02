package config

import (
	_ "embed"
	"os"
	"sync"
)

//go:embed data/qqwry.dat
var ipv4QQWry []byte

//go:embed data/ipv6wry.db
var ipv6QQWry []byte

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

// GetConfigQQWry
// @autowire(set=config)
func GetConfigQQWry() QQWry {
	result := GetConfig().QQWry
	result.IPv4Data = ipv4QQWry
	result.IPv6Data = ipv6QQWry
	return result
}

// GetConfigBackendList
// @autowire(set=config)
func GetConfigBackendList() BackendList {
	list := GetConfig().BackendList
	list.Default()
	return list
}
