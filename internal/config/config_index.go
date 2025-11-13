package config

import "github.com/Lofanmi/gobana/service"

// ---------------------------------------------------------------------------------------------------------------------

type IndexName = string

// ---------------------------------------------------------------------------------------------------------------------

type IndexConfig struct {
	Label         string               `json:"label"`          // 用于界面展示
	Name          IndexName            `json:"name"`           // 日志存储名称，如 access-log/json-log/string-log 等
	Meta          IndexMeta            `json:"meta"`           // 额外配置
	ProviderName  service.ProviderName `json:"provider_name"`  // 关联的查询器
	ParserName    ParserName           `json:"parser_name"`    // 关联的解析器
	BuildInQuery  BuildInQuery         `json:"build_in_query"` // 内置的快捷查询
	DefaultFields []string             `json:"default_fields"` // 默认查询字段
}

type Indexes map[IndexName]IndexConfig

func (s *Indexes) Register(c *IndexConfig) {
	(*s)[c.Name] = *c
}

// ---------------------------------------------------------------------------------------------------------------------

type IndexMeta struct {
	IndexMetaForSls     *IndexMetaForSls     `json:"index_meta_for_sls,omitempty"`
	IndexMetaForElastic *IndexMetaForElastic `json:"index_meta_for_elastic,omitempty"`
}

type IndexMetaForSls struct {
	Project string `json:"project"`
	Store   string `json:"store"`
}

type IndexMetaForElastic struct {
	Indexes    []string           `json:"indexes"`
	TimeField  string             `json:"time_field"`  // 时间排序字段（@timestamp）
	Timezone   string             `json:"timezone"`    // 时区
	SortFields []ElasticSortField `json:"sort_fields"` // 排序字段配置
}

// ---------------------------------------------------------------------------------------------------------------------

type ElasticSortField struct {
	Field     string `yaml:"field"`
	Ascending bool   `yaml:"ascending"`
}

// ---------------------------------------------------------------------------------------------------------------------
