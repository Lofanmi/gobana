package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lofanmi/gobana/internal/constant"
	"github.com/spf13/cast"
)

type Config struct {
	Application Application `yaml:"application"`
	BackendList BackendList `yaml:"backend_list"`
}

type Application struct {
	Production bool   `yaml:"production"`
	ListenAddr string `yaml:"listen_addr"`
}

type BackendList []Backend

type Backend struct {
	Name           string                  `yaml:"name"`
	Enabled        bool                    `yaml:"enabled"`
	Type           string                  `yaml:"type"`
	Addr           string                  `yaml:"addr"`
	Auth           Auth                    `yaml:"auth"`
	Timeout        int64                   `yaml:"timeout"`          // 请求超时时间（毫秒）
	MultiSearch    map[string]MultiSearch  `yaml:"multi_search"`     // 多索引/日志存储查询
	BuildInQueries map[string]BuildInQuery `yaml:"build_in_queries"` // 内置的快捷查询
	TimeField      map[string]string       `yaml:"time_field"`       // 时间排序字段，ES默认为@timestamp。
	Timezone       map[string]string       `yaml:"timezone"`         // 时区
	DefaultFields  map[string][]string     `yaml:"default_fields"`   // 默认查询字段
	SortFields     map[string][]SortField  `yaml:"sort_fields"`      // 字段排序
	ParserLogType  string                  `yaml:"parser_log_type"`  // 日志类型解析器
	ParserFields   ParserFields            `yaml:"parser_fields"`    // 字段解析器
}

type Auth struct {
	Enabled         bool   `yaml:"enabled"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	KbnVersion      string `yaml:"kbn_version"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Cookie          string `yaml:"cookie"`
}

type MultiSearch struct {
	Name      string   `yaml:"name"`
	Order     int      `yaml:"order"`
	IndexList []string `yaml:"index_list"`
}

type BuildInQuery struct {
	Must    BuildInQueryEntrySlice `yaml:"must"`
	MustNot BuildInQueryEntrySlice `yaml:"must_not"`
	Or      BuildInQueryEntrySlice `yaml:"or"`
}

type BuildInQueryEntrySlice []BuildInQueryEntry

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

type ParserField struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`
	FromFields []string `yaml:"from_field"`
	ToField    string   `yaml:"to_field"`
	TrimSet    string   `yaml:"trim_set"`
	LuaField   string   `yaml:"lua_field"`
	LuaReturn  string   `yaml:"lua_return"`
}

type kibanaAuth struct {
	ProviderType string `json:"providerType"`
	ProviderName string `json:"providerName"`
	CurrentURL   string `json:"currentURL"`
	Params       struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"params"`
}

type MultiSearchSlice []MultiSearch

func (s MultiSearchSlice) Len() int           { return len(s) }
func (s MultiSearchSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MultiSearchSlice) Less(i, j int) bool { return s[i].Order < s[j].Order }

func (s *Backend) Default() {
	if s.Timeout <= 0 {
		s.Timeout = 60 * 1000
	}
	s.ParserFields.Default()
	s.authIfNeeded()
}

func (s *Backend) authIfNeeded() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	if !s.Auth.Enabled {
		return
	}
	if s.Type != constant.ClientTypeKibanaProxy {
		return
	}
	if s.Auth.Cookie != "" {
		return
	}
	var auth kibanaAuth
	auth.ProviderType = "basic"
	auth.ProviderName = "basic"
	auth.CurrentURL = s.Addr + "/login?next=%2Fapp%2Fdiscover#/"
	auth.Params.Username = s.Auth.Username
	auth.Params.Password = s.Auth.Password
	var data []byte
	if data, err = json.Marshal(&auth); err != nil {
		return
	}
	var request *http.Request
	if request, err = http.NewRequest(http.MethodPost, s.Addr+"/internal/security/login", bytes.NewReader(data)); err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("kbn-version", s.Auth.KbnVersion)
	var response *http.Response
	if response, err = http.DefaultClient.Do(request); err != nil {
		return
	}
	defer func() { _ = response.Body.Close() }()
	for _, cookie := range response.Cookies() {
		s.Auth.Cookie = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
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

func (s *ParserFields) Default() {
	for _, field := range s.AccessLog {
		field.Default()
	}
	for _, field := range s.JsonLog {
		field.Default()
	}
	for _, field := range s.StringLog {
		field.Default()
	}
}

func (s *ParserField) Default() {
}

func (s *BuildInQueryEntrySlice) FieldConditions(fn func(field string, conditions []string) bool) {
	for _, item := range *s {
		if !item.Always {
			continue
		}
		if !fn(item.Field, cast.ToStringSlice(item.Values)) {
			break
		}
	}
}
