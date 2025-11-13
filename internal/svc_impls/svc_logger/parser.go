package svc_logger

import (
	"strings"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/dop251/goja"
	"github.com/spf13/cast"
	"github.com/tidwall/sjson"
)

func (s *Service) doParseSLS(indexName string, r service.SearchResultSls) (logs service.LogItems, err error) {
	if r.ResponseLog == nil || len(r.ResponseLog.Logs) <= 0 {
		return
	}
	index := s.Indexes[indexName]
	for _, _hit := range r.ResponseLog.Logs {
		data := slsLogToMapStringAny(_hit)
		data["_index"] = indexName
		logItem := service.LogItem{Storage: index.Label}
		if err = s.parseLogEntrySls(indexName, data, &logItem); err != nil {
			return
		}
		logs = append(logs, logItem)
	}
	return
}

func (s *Service) parseLogEntrySls(indexName string, data map[string]any, logItem *service.LogItem) (err error) {
	*logItem = service.LogItem{
		Timestamp: 0,
		Log:       data,
	}
	return
}

// func (s *Service) doParseElastic(indexName string, r *elastic.SearchResult) (logs service.LogItems, err error) {
// 	if r == nil || r.Hits == nil {
// 		return
// 	}
// 	index := s.Indexes[indexName]
// 	for _, hit := range r.Hits.Hits {
// 		logItem := service.LogItem{Storage: index.Label}
// 		if err = s.parseLogEntryElastic(indexName, data, &logItem); err != nil {
// 			return
// 		}
// 		logs = append(logs, logItem)
// 	}
// 	return
// }

func (s *Service) parseFieldsFromGoJaObject(field *config.ParserField, data map[string]any, targetJSON *string, _sourceObject *goja.Object) {
	switch field.Type {
	case service.ParserFieldTypeReplacements:
		var value string
		for _, fromField := range field.FromFields {
			value = cast.ToString(data[fromField])
			if field.TrimSet != "" {
				value = strings.Trim(value, field.TrimSet)
			}
			if value != "" {
				break
			}
		}
		if newValue, err := sjson.Set(*targetJSON, field.ToField, value); err == nil {
			*targetJSON = newValue
		}
	case service.ParserFieldTypeJavaScript:
		vm := s.GoJa.GetRuntime()
		for _, _fromField := range field.FromFields {
			fromField := cast.ToString(data[_fromField])
			if fromField == "" {
				continue
			}
			program, err := goja.Compile("", field.JavaScriptField, false)
			if err != nil {
				continue
			}
			_ = vm.Set("value", fromField)
			_ = vm.Set("source", _sourceObject)
			ret, err := vm.RunProgram(program)
			if err != nil {
				continue
			}
			var newValue any
			switch field.JavaScriptReturn {
			case service.ParserFieldReturnString:
				if str, ok := ret.Export().(string); ok {
					newValue = str
				} else {
					continue
				}
			case service.ParserFieldReturnNumber:
				if num, ok := ret.Export().(float64); ok {
					newValue = num
				} else {
					continue
				}
			default:
				continue
			}
			if newJSON, err := sjson.Set(*targetJSON, field.ToField, newValue); err == nil {
				*targetJSON = newJSON
			}
		}
	}
}

func slsLogToMapStringAny(m map[string]string) (res map[string]any) {
	res = map[string]any{}
	for k, v := range m {
		if v != "" && v != "null" {
			res[k] = v
		}
	}
	return
}
