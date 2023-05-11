package config

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	lua "github.com/yuin/gopher-lua"
)

type Config struct {
	BackendList BackendList `yaml:"backend_list"`
}

type BackendList []Backend

type Backend struct {
	Name           string                  `yaml:"name"`
	Type           string                  `yaml:"type"`
	Addr           string                  `yaml:"addr"`
	Auth           Auth                    `yaml:"auth"`
	Timeout        int64                   `yaml:"timeout"`          // 请求超时时间（毫秒）
	MultiSearch    map[string]MultiSearch  `yaml:"multi_search"`     // 多索引/日志存储查询
	BuildInQueries map[string]BuildInQuery `yaml:"build_in_queries"` // 内置的快捷查询
	TimeField      map[string]string       `yaml:"time_field"`       // 时间排序字段
	DefaultFields  map[string][]string     `yaml:"default_fields"`   // 默认查询字段
	SortFields     map[string][]SortField  `yaml:"sort_fields"`      // 字段排序
	ParserFields   ParserFields            `yaml:"parser_fields"`    // 字段解析器
}

type Auth struct {
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

type MultiSearch struct {
	IndexList    []string `yaml:"index_list"`
	Project      string   `yaml:"project"`
	LogStoreList []string `yaml:"log_store_list"`
}

type BuildInQuery struct {
	Must    []BuildInQueryEntry `yaml:"must"`
	MustNot []BuildInQueryEntry `yaml:"must_not"`
	Or      []BuildInQueryEntry `yaml:"or"`
}

type BuildInQueryEntry struct {
	Name     string        `yaml:"name"`
	Field    string        `yaml:"field"`
	Values   []interface{} `yaml:"values"`
	Operator string        `yaml:"operator"`
	Always   bool          `yaml:"always"`
}

type SortField struct {
	Field     string `yaml:"field"`
	Ascending bool   `yaml:"ascending"`
}

type ParserFields struct {
	AccessLog []ParserField `yaml:"access_log"`
	JsonLog   []ParserField `yaml:"json_log"`
	StringLog []ParserField `yaml:"string_log"`
}

type (
	ParserFieldType   = string
	ParserFieldReturn = string
)

const (
	ParserFieldTypeReplacements ParserFieldType = "replacements"
	ParserFieldTypeLua          ParserFieldType = "lua"

	ParserFieldReturnString ParserFieldReturn = "string"
	ParserFieldReturnNumber ParserFieldReturn = "number"
)

type ParserField struct {
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Return     string      `yaml:"return"`
	FromFields []string    `yaml:"from_field"`
	ToField    string      `yaml:"to_field"`
	Lua        string      `yaml:"lua"`
	L          *lua.LState `yaml:"-"`
}

func (s *Backend) Default(L *lua.LState) {
	if s.Timeout <= 0 {
		s.Timeout = 60 * 1000
	}
	s.ParserFields.Default(L)
}

func (s *BackendList) Default(L *lua.LState) {
	if len(*s) <= 0 {
		return
	}
	for i := 0; i < len(*s); i++ {
		(*s)[i].Default(L)
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

func (s *ParserFields) Default(L *lua.LState) {
	for _, field := range s.AccessLog {
		field.Default(L)
	}
	for _, field := range s.JsonLog {
		field.Default(L)
	}
	for _, field := range s.StringLog {
		field.Default(L)
	}
}

func (s *ParserField) Default(L *lua.LState) {
	s.L = L
}

func (s *ParserField) Handle(g gjson.Result, targetJSON *string, source string) {
	switch s.Type {
	case ParserFieldTypeReplacements:
		for _, fromField := range s.FromFields {
			value := g.Get(fromField).String()
			if value == "" {
				continue
			}
			if newValue, err := sjson.Set(*targetJSON, s.ToField, value); err == nil {
				*targetJSON = newValue
			}
		}
	case ParserFieldTypeLua:
		L := s.L
		for _, fromField := range s.FromFields {
			value := g.Get(fromField).String()
			if value == "" {
				continue
			}
			if err := L.CallByParam(lua.P{Fn: L.GetGlobal("parse_field"), NRet: 2, Protect: true}, lua.LString(value), lua.LString(source)); err != nil {
				continue
			}
			ret, errString := L.Get(-2), L.Get(-1)
			L.Pop(2)
			if errString.String() != "" {
				continue
			}
			var newValue interface{}
			switch s.Return {
			case ParserFieldReturnString:
				if res, ok := ret.(lua.LString); !ok {
					continue
				} else {
					newValue = string(res)
				}
			case ParserFieldReturnNumber:
				if res, ok := ret.(lua.LNumber); !ok {
					continue
				} else {
					newValue = float64(res)
				}
			default:
				continue
			}
			if newJSON, err := sjson.Set(*targetJSON, s.ToField, newValue); err == nil {
				*targetJSON = newJSON
			}
		}
	}
}
