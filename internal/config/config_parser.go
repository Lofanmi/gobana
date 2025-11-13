package config

import (
	"github.com/Lofanmi/gobana/service"
)

// ---------------------------------------------------------------------------------------------------------------------

type ParserConfig struct {
	Name     service.ParserName `json:"name"`     // 解析器名称
	BaseDir  string             `json:"base_dir"` // 字段解析器脚本路径
	Filename string             `json:"filename"` // 日志解析器
	Fields   []ParserField      `json:"fields"`   // 字段列表
}

type ParserField struct {
	Name             string                  `json:"name"`
	Type             service.ParserFieldType `json:"type"`
	FromFields       []string                `json:"from_field"`
	ToField          string                  `json:"to_field"`
	TrimSet          string                  `json:"trim_set"`
	JavaScriptField  string                  `json:"javascript_field"` // 记录脚本函数，用于字段处理
	JavaScriptReturn string                  `json:"javascript_return"`
}

// ---------------------------------------------------------------------------------------------------------------------
