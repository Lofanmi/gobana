package svc_logger

import (
	"errors"
	"strings"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/dop251/goja"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

func (s *Service) doParseSLS(indexName string, r service.SearchResultSls) (logs service.LogItems, err error) {
	if r.ResponseLog == nil || len(r.ResponseLog.Logs) <= 0 {
		return
	}
	for _, _hit := range r.ResponseLog.Logs {
		data := slsLogToMapStringAny(_hit)
		data["_index"] = indexName
		logItem, e := s.parseLogEntry(indexName, data)
		if e != nil {
			err = e
			return
		}
		logs = append(logs, logItem)
	}
	return
}

func (s *Service) doParseElastic(indexName string, r service.SearchResultElastic) (logs service.LogItems, err error) {
	if r.Response == nil || r.Response.Hits == nil {
		return
	}
	for _, _hit := range r.Response.Hits.Hits {
		var data map[string]any
		if err = jsoniter.Unmarshal(_hit.Source, &data); err != nil {
			return
		}
		data["_index"] = indexName
		logItem, e := s.parseLogEntry(indexName, data)
		if e != nil {
			err = e
			return
		}
		logs = append(logs, logItem)
	}
	return
}

func (s *Service) parseLogEntry(indexName string, data map[string]any) (logItem service.LogItem, err error) {
	index := s.Indexes[indexName]
	parserConfig, ok := s.Parsers[index.ParserName]
	if !ok || len(parserConfig.Fields) <= 0 {
		err = errors.New("no parser found for " + indexName)
		return
	}
	logItem.Storage = index.Label
	logItem.Source = data
	logItem.Log = make(map[string]any)
	for _, field := range parserConfig.Fields {
		err = s.parseFieldsFromGoJaObject(parserConfig.Name, field, &logItem)
		if err != nil {
			return
		}
	}
	logItem.Timestamp = cast.ToInt64(logItem.Log[service.GobanaTimestampField])
	return
}

func (s *Service) parseFieldsFromGoJaObject(parserName string, field config.ParserField, logItem *service.LogItem) (err error) {
	switch field.Type {
	case service.ParserFieldTypeTimestamp:
		vm, callable, e := s.GoJa.GetCallable(parserName, service.GobanaTimestampFunc)
		if e != nil {
			err = e
			return
		}
		result, e := callable(goja.Undefined(), vm.ToValue(logItem.Source))
		if e != nil {
			err = e
			return
		}
		ts := cast.ToFloat64(result.Export())
		logItem.Log[service.GobanaTimestampField] = ts
		logItem.Timestamp = cast.ToInt64(ts)
	case service.ParserFieldTypeReplacements:
		var value string
		for _, fromField := range field.FromFields {
			value = cast.ToString(logItem.Source[fromField])
			if field.TrimSet != "" {
				value = strings.Trim(value, field.TrimSet)
			}
			if value != "" {
				break
			}
		}
		logItem.Log[field.ToField] = value
	case service.ParserFieldTypeJavaScript:
		vm, callable, e := s.GoJa.GetCallable(parserName, field.JavaScriptField)
		if e != nil {
			err = e
			return
		}
		result, e := callable(goja.Undefined(), vm.ToValue(logItem.Source))
		if e != nil {
			err = e
			return
		}
		var newValue any
		exportedResult := result.Export()
		switch field.JavaScriptReturn {
		case service.ParserFieldReturnString:
			newValue = cast.ToString(exportedResult)
		case service.ParserFieldReturnNumber:
			newValue = cast.ToFloat64(exportedResult)
		}
		logItem.Log[field.ToField] = newValue
	}
	return
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
