package config

import (
	"os"

	"github.com/Lofanmi/gobana/service"
	"github.com/dop251/goja"
)

// ---------------------------------------------------------------------------------------------------------------------

type ParserConfig struct {
	Name               service.ParserName `json:"name"`                // 解析器名称
	JavaScriptFilename string             `json:"javascript_filename"` // JavaScript文件名，默认main.js
	Fields             []ParserField      `json:"fields"`              // 字段列表
	Program            *goja.Program      `json:"-"`                   // 私有字段，存储编译后的程序
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

type ParserName = string

type Parsers map[ParserName]ParserConfig

func (c *ParserConfig) init() {
	if c.JavaScriptFilename == "" {
		panic("JavaScriptFilename is empty")
	}
	content, _ := os.ReadFile(c.JavaScriptFilename)
	if len(content) <= 0 {
		panic("JavaScriptFile is empty")
	}
	program, err := goja.Compile(c.JavaScriptFilename, string(content), false)
	if err != nil {
		panic(err)
	}
	c.Program = program
}

// ---------------------------------------------------------------------------------------------------------------------
