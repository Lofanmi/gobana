package config

// ---------------------------------------------------------------------------------------------------------------------

type IndexName = string

// ---------------------------------------------------------------------------------------------------------------------

type IndexConfig struct {
	Label        string       `json:"label"`         // 用于界面展示
	Name         IndexName    `json:"name"`          // 日志存储名称，如 access-log/json-log/string-log 等
	Meta         IndexMeta    `json:"meta"`          // 额外配置
	ProviderName ProviderName `json:"provider_name"` // 关联的查询器
	ParserName   ParserName   `json:"parser_name"`   // 关联的解析器
}

type Indexes map[IndexName]IndexConfig

func (s *Indexes) Register(c *IndexConfig) {
	(*s)[c.Name] = *c
}

// ---------------------------------------------------------------------------------------------------------------------

type IndexMeta struct {
	IndexMetaForSls *IndexMetaForSls `json:"index_meta_for_sls,omitempty"`
}

type IndexMetaForSls struct {
	Project string `json:"project"`
	Store   string `json:"store"`
}

// ---------------------------------------------------------------------------------------------------------------------
