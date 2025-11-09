package config

import (
	"github.com/Lofanmi/gobana/service"
)

// ---------------------------------------------------------------------------------------------------------------------

type Providers map[service.ProviderName]ProviderConfig

func (s Providers) Match(name string) (res ProviderConfig) {
	res, _ = s[name]
	return
}

func (s Providers) Default() {
	if len(s) <= 0 {
		return
	}
	for _, c := range s {
		c.Default()
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type ProviderConfig struct {
	Label              string               `json:"label"`                 // 用于前端展示
	Name               service.ProviderName `json:"name"`                  // 名称
	Type               service.ProviderType `json:"type"`                  // 后端类型
	TimeoutMillisecond int64                `json:"timeout,omitempty"`     // 请求超时时间（毫秒）
	AuthConfig         AuthConfig           `json:"auth_config,omitempty"` // 后端认证配置
}

func (s *ProviderConfig) Default() {
	if s.TimeoutMillisecond <= 0 {
		s.TimeoutMillisecond = 60 * 1000
	}
	err := s.authIfNeeded()
	if err != nil {
		panic(err)
	}
}

func (s *ProviderConfig) authIfNeeded() (err error) {
	switch s.Type {
	case service.ProviderTypeKibanaProxy:
		err = s.AuthConfig.AuthForKibanaProxy.init()
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
