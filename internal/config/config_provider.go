package config

// ---------------------------------------------------------------------------------------------------------------------

type ProviderName = string

// ---------------------------------------------------------------------------------------------------------------------

type ProviderType = string

const (
	ProviderTypeSls           ProviderType = "sls"
	ProviderTypeSlsProxy      ProviderType = "sls_proxy"
	ProviderTypeElasticsearch ProviderType = "elasticsearch"
	ProviderTypeKibanaProxy   ProviderType = "kibana_proxy"
)

// ---------------------------------------------------------------------------------------------------------------------

type Providers map[ProviderName]ProviderConfig

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
	Label              string       `json:"label"`                 // 用于前端展示
	Name               ProviderName `json:"name"`                  // 名称
	Type               ProviderType `json:"type"`                  // 后端类型
	TimeoutMillisecond int64        `json:"timeout,omitempty"`     // 请求超时时间（毫秒）
	AuthConfig         AuthConfig   `json:"auth_config,omitempty"` // 后端认证配置
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
	case ProviderTypeKibanaProxy:
		err = s.AuthConfig.AuthForKibanaProxy.init()
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
