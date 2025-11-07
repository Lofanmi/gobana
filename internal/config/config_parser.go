package config

type ParserName = string

// ---------------------------------------------------------------------------------------------------------------------

type ParserType = string

const (
	ParserTypeReplacements ParserType = "replacements"
	ParserTypeJavascript   ParserType = "javascript"
)

// ---------------------------------------------------------------------------------------------------------------------

type ParserConfig struct {
	Name     ParserName    `json:"name"`     // 解析器名称
	BaseDir  string        `json:"base_dir"` // 字段解析器脚本路径
	Filename string        `json:"filename"` // 日志解析器
	Fields   []ParserField `json:"fields"`   // 字段列表
}

type ParserField struct {
	Name             string     `json:"name"`
	Type             ParserType `json:"type"`
	FromFields       []string   `json:"from_field"`
	ToField          string     `json:"to_field"`
	TrimSet          string     `json:"trim_set"`
	JavaScriptField  string     `json:"javascript_field"`
	JavaScriptReturn string     `json:"javascript_return"`
}

// ---------------------------------------------------------------------------------------------------------------------
