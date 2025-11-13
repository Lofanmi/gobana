package config

// ---------------------------------------------------------------------------------------------------------------------

type BackendName = string

// ---------------------------------------------------------------------------------------------------------------------

type BackendConfig struct {
	Label   string      `json:"label"` // eg. 国内SLS-访问日志/程序日志
	Name    BackendName `json:"name"`  // eg. cn_sls_access_log_json_log
	Order   int         `json:"order"`
	Indexes []IndexName `json:"indexes"`
}

type Backends map[BackendName]BackendConfig

func (s *Backends) Register(c *BackendConfig) {
	(*s)[c.Name] = *c
}

// ---------------------------------------------------------------------------------------------------------------------

type BackendSlice []BackendConfig

func (s BackendSlice) Len() int           { return len(s) }
func (s BackendSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BackendSlice) Less(i, j int) bool { return s[i].Order < s[j].Order }

// ---------------------------------------------------------------------------------------------------------------------
