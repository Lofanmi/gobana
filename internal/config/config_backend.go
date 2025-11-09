package config

// ---------------------------------------------------------------------------------------------------------------------

type BackendName = string

// ---------------------------------------------------------------------------------------------------------------------

type BackendConfig struct {
	Label         string                     `json:"label"` // eg. 国内SLS-访问日志/程序日志
	Name          BackendName                `json:"name"`  // eg. cn_sls_access_log_json_log
	Order         int                        `json:"order"`
	Indexes       []IndexName                `json:"indexes"`
	TimeField     map[IndexName]string       `json:"time_field"`     // 时间排序字段，ES默认为@timestamp。
	Timezone      map[IndexName]string       `json:"timezone"`       // 时区
	DefaultFields map[IndexName][]string     `json:"default_fields"` // 默认查询字段
	BuildInQuery  map[IndexName]BuildInQuery `json:"build_in_query"` // 内置的快捷查询
	SortFields    map[IndexName][]SortField  `json:"sort_fields"`    // 字段排序

	// MultiSearch    map[string]MultiSearch  `yaml:"multi_search"`     // 多索引/日志存储查询
	// ParserLogType  string                     `yaml:"parser_log_type"`  // 日志类型解析器
	// ParserFields   ParserFields               `yaml:"parser_fields"`    // 字段解析器
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

type SortField struct {
	Field     string `yaml:"field"`
	Ascending bool   `yaml:"ascending"`
}

// ---------------------------------------------------------------------------------------------------------------------
