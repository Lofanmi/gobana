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
		if err := DefaultLoader.Load(os.Getenv("CONFIG"), &config); err != nil {
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

// GetConfigProviders
// @autowire(set=config)
func GetConfigProviders() Providers { return GetConfig().Providers }

// GetConfigBackends
// @autowire(set=config)
func GetConfigBackends() Backends { return GetConfig().Backends }

// GetConfigIndexes
// @autowire(set=config)
func GetConfigIndexes() Indexes { return GetConfig().Indexes }
