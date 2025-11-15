package config

import (
	"os"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

var (
	DefaultLoader Loader = &localFileLoader{}
)

type Loader interface {
	Load(filename string, dst *Config) (err error)
}

type localFileLoader struct{}

func (localFileLoader) Load(filename string, dst *Config) (err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	processedData, err := processJSON(data)
	if err != nil {
		return
	}
	return jsoniter.Unmarshal(processedData, dst)
}

func processJSON(data []byte) (res []byte, err error) {
	parsed := gjson.ParseBytes(data)
	result := processValue(parsed)
	return []byte(result.Raw), nil
}

func processValue(value gjson.Result) gjson.Result {
	switch value.Type {
	case gjson.String:
		resolved := resolveEnvVars(value.String())
		if resolved == "true" {
			return gjson.Parse("true")
		}
		if resolved == "false" {
			return gjson.Parse("false")
		}
		return gjson.Parse(`"` + resolved + `"`)
	case gjson.JSON:
		if value.IsObject() {
			var objBuilder strings.Builder
			objBuilder.WriteString("{")
			first := true
			value.ForEach(func(key, val gjson.Result) bool {
				if !first {
					objBuilder.WriteString(",")
				}
				first = false
				objBuilder.WriteString(`"` + key.String() + `":`)
				objBuilder.WriteString(processValue(val).Raw)
				return true
			})
			objBuilder.WriteString("}")
			return gjson.Parse(objBuilder.String())
		} else if value.IsArray() {
			var arrayBuilder strings.Builder
			arrayBuilder.WriteString("[")
			first := true
			value.ForEach(func(_, val gjson.Result) bool {
				if !first {
					arrayBuilder.WriteString(",")
				}
				first = false
				arrayBuilder.WriteString(processValue(val).Raw)
				return true
			})
			arrayBuilder.WriteString("]")
			return gjson.Parse(arrayBuilder.String())
		}
	default:
		return value
	}
	return value
}

func resolveEnvVars(input string) string {
	re := regexp.MustCompile(`\$\{env\.([^}]+)}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		envVarName := match[6 : len(match)-1]
		if envValue := os.Getenv(envVarName); envValue != "" {
			return envValue
		}
		return match
	})
}
