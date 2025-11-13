package svc_logger

import (
	"encoding/json"
	"strings"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/internal/gotil/lua_json"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	lua "github.com/yuin/gopher-lua"
)

func (s *Service) doParseSLS(indexName string, r service.SearchResultSls) (logs service.LogItems, err error) {
	if r.ResponseLog == nil || len(r.ResponseLog.Logs) <= 0 {
		return
	}
	index := s.Indexes[indexName]
	for _, hit := range r.ResponseLog.Logs {
		for k, v := range hit {
			if v == "" || v == "null" {
				delete(hit, k)
			}
		}
		hit["_index"] = indexName
		data, _ := json.Marshal(hit)
		logItem := service.LogItem{Storage: index.Label}
		if err = s.parseLogBytesSls(indexName, data, &logItem); err != nil {
			return
		}
		logs = append(logs, logItem)
	}
	return
}

func (s *Service) doParseElastic(indexName string, r *elastic.SearchResult) (logs service.LogItems, err error) {
	if r == nil || r.Hits == nil {
		return
	}
	index := s.Indexes[indexName]
	for _, hit := range r.Hits.Hits {
		data, _ := json.Marshal(hit)
		logItem := service.LogItem{Storage: index.Label}
		if err = s.parseLogBytesElastic(indexName, data, &logItem); err != nil {
			return
		}
		logs = append(logs, logItem)
	}
	return
}

func (s *Service) parseLogBytesSls(indexName string, data []byte, logItem *service.LogItem) (err error) {
	hitMap := map[string]any{}
	if err = json.Unmarshal(data, &hitMap); err != nil {
		return
	}
	tb := gotil.MapToTable(hitMap)
	var logInterface any
	if _sourceString != "" {
		data = []byte(_sourceString)
		if err = json.Unmarshal(data, &hitMap); err != nil {
			return
		}
	}
	logItem = service.LogItem{
		Timestamp: gotil.ParseTime(logTime),
		LogType:   logType,
		Log:       logInterface,
	}
	return
}

func (s *Service) handleParserField(field *config.ParserField, g gjson.Result, targetJSON *string, _sourceTable *lua.LTable) {
	switch field.Type {
	case service.ParserFieldTypeReplacements:
		var value string
		for _, fromField := range field.FromFields {
			value = g.Get(fromField).String()
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
	case service.ParserFieldTypeLua:
		L, fn := s.GetLuaState()
		lua_json.Preload(L)
		defer fn()
		for _, fromField := range field.FromFields {
			value := g.Get(fromField).String()
			if value == "" {
				continue
			}
			if err := L.DoString(field.LuaField); err != nil {
				continue
			}
			if err := L.CallByParam(lua.P{Fn: L.GetGlobal("parse_field"), NRet: 2, Protect: true}, lua.LString(value), _sourceTable); err != nil {
				continue
			}
			ret, errString := L.Get(-2), L.Get(-1)
			L.Pop(2)
			if errString.String() != "" {
				continue
			}
			var newValue any
			switch field.LuaReturn {
			case service.ParserFieldReturnString:
				if res, ok := ret.(lua.LString); !ok {
					continue
				} else {
					newValue = string(res)
				}
			case service.ParserFieldReturnNumber:
				if res, ok := ret.(lua.LNumber); !ok {
					continue
				} else {
					newValue = float64(res)
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
