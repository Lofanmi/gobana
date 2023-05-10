package config

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Config struct {
	BackendList BackendList `json:"backend_list"`
}

type BackendList []Backend

type Backend struct {
	Name           string                  `json:"name"`
	Type           string                  `json:"type"`
	Addr           string                  `json:"addr"`
	Auth           Auth                    `json:"auth"`
	Timeout        int64                   `json:"timeout"`          // 请求超时时间（毫秒）
	MultiSearch    map[string]MultiSearch  `json:"multi_search"`     // 多索引/日志存储查询
	BuildInQueries map[string]BuildInQuery `json:"build_in_queries"` // 内置的快捷查询
	TimeField      map[string]string       `json:"time_field"`       // 时间排序字段
	DefaultFields  map[string][]string     `json:"default_fields"`   // 默认查询字段
	SortFields     map[string][]SortField  `json:"sort_fields"`      // 字段排序
	ParserFields   map[string][]string
}

type Auth struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type MultiSearch struct {
	IndexList    []string `json:"index_list"`
	Project      string   `json:"project"`
	LogStoreList []string `json:"log_store_list"`
}

type BuildInQuery struct {
	Must    []BuildInQueryEntry `json:"must"`
	MustNot []BuildInQueryEntry `json:"must_not"`
	Or      []BuildInQueryEntry `json:"or"`
}

type BuildInQueryEntry struct {
	Name     string        `json:"name"`
	Field    string        `json:"field"`
	Values   []interface{} `json:"values"`
	Operator string        `json:"operator"`
	Always   bool          `json:"always"`
}

type SortField struct {
	Field     string `json:"field"`
	Ascending bool   `json:"ascending"`
}

type ParserFieldType = string

const (
	ParserFieldTypeReplacement ParserFieldType = "replacement"
	ParserFieldTypeScript      ParserFieldType = "script"
)

type ParserField struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	FromField string `json:"from_field"`
	ToField   string `json:"to_field"`
	Script    string `json:"script"`
}

func (s *Backend) Default() {
	if s.Timeout <= 0 {
		s.Timeout = 60 * 1000
	}
}

func (s *BackendList) Default() {
	if len(*s) <= 0 {
		return
	}
	for i := 0; i < len(*s); i++ {
		(*s)[i].Default()
	}
}

func (s *BackendList) Match(name string) Backend {
	if len(*s) <= 0 {
		return Backend{}
	}
	for _, backend := range *s {
		if backend.Name == name {
			return backend
		}
	}
	return Backend{}
}

func (s *ParserField) Handle(g gjson.Result, targetJSON *string) {
	switch s.Type {
	case ParserFieldTypeReplacement:
		value := g.Get(s.FromField).String()
		if newValue, err := sjson.Set(*targetJSON, s.ToField, value); err == nil {
			*targetJSON = newValue
		}
	case ParserFieldTypeScript:
	}
}
