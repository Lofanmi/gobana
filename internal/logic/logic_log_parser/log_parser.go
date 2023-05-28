package logic_log_parser

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"unsafe"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/constant"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/internal/gotil/lua_json"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	lua "github.com/yuin/gopher-lua"
)

var (
	_ logic.LogParser = &LogParser{}
)

// LogParser
// @autowire(logic.LogParser,set=logics)
type LogParser struct {
	LuaState logic.LuaState
}

func (s *LogParser) ParseElastic(backend config.Backend, m map[string]*elastic.SearchResult) (total int, logs service.LogItems, err error) {
	total = s.parseElasticTotal(m)
	logs = make([]service.LogItem, 0, total)
	for _, result := range m {
		if result == nil || result.Hits == nil {
			continue
		}
		for _, hit := range result.Hits.Hits {
			data, _ := json.Marshal(hit)
			var logItem service.LogItem
			if logItem, err = s.parseLogBytes(backend, data); err != nil {
				return
			}
			logItem.Storage = hit.Index
			logs = append(logs, logItem)
		}
	}
	sort.Sort(logs)
	return
}

func (s *LogParser) ParseSLS(backend config.Backend, m map[string]logic.SLSSearchResult) (total int, logs service.LogItems, err error) {
	logs = make([]service.LogItem, 0, total)
	for index, result := range m {
		countRes, logRes := result.ResponseCount, result.ResponseLog
		if logRes == nil || len(logRes.Logs) <= 0 {
			continue
		}
		if countRes != nil && len(countRes.Logs) > 0 {
			log := countRes.Logs[0]
			total = cast.ToInt(log["count"])
		} else {
			total = 10000
		}
		for _, hit := range logRes.Logs {
			for k, v := range hit {
				if v == "" || v == "null" {
					delete(hit, k)
				}
			}
			hit["_index"] = index
			data, _ := json.Marshal(hit)
			var logItem service.LogItem
			if logItem, err = s.parseLogBytes(backend, data); err != nil {
				return
			}
			logItem.Storage = index
			logs = append(logs, logItem)
		}
	}
	sort.Sort(logs)
	return
}

func (s *LogParser) parseElasticTotal(m map[string]*elastic.SearchResult) (total int) {
	for _, result := range m {
		if result == nil || result.Hits == nil || result.Hits.TotalHits == nil {
			continue
		}
		total += int(result.Hits.TotalHits.Value)
	}
	return
}

func (s *LogParser) parseLogBytes(backend config.Backend, data []byte) (logItem service.LogItem, err error) {
	hitMap := map[string]interface{}{}
	if err = json.Unmarshal(data, &hitMap); err != nil {
		return
	}
	tb := gotil.MapToTable(hitMap)
	var logInterface interface{}
	logTime := ""
	logType, _sourceTable, _sourceString, e := s.parseLogType(backend, tb)
	if e != nil {
		err = e
		return
	}
	if _sourceString != "" {
		data = []byte(_sourceString)
		if err = json.Unmarshal(data, &hitMap); err != nil {
			return
		}
	}
	switch logType {
	case service.LogTypeAccessLog:
		log := new(service.AccessLog)
		_ = parseLog(s, backend.ParserFields.AccessLog, data, hitMap, _sourceTable, log)
		logInterface, logTime = log, log.Time
	case service.LogTypeJsonLog:
		log := new(service.JsonLog)
		_ = parseLog(s, backend.ParserFields.JsonLog, data, hitMap, _sourceTable, log)
		logInterface, logTime = log, log.Time
	case service.LogTypeStringLog:
		log := new(service.StringLog)
		_ = parseLog(s, backend.ParserFields.StringLog, data, hitMap, _sourceTable, log)
		logInterface, logTime = log, log.Time
	default:
		_ = e
		return
	}
	logItem = service.LogItem{
		Timestamp: gotil.ParseTime(logTime),
		LogType:   logType,
		Log:       logInterface,
	}
	return
}

func (s *LogParser) parseLogType(backend config.Backend, tb *lua.LTable) (logType service.LogType, _sourceTable *lua.LTable, _sourceString string, err error) {
	L, fn := s.LuaState.GetLuaState()
	lua_json.Preload(L)
	defer fn()
	if err = L.DoString(backend.ParserLogType); err != nil {
		return
	}
	if err = L.CallByParam(lua.P{Fn: L.GetGlobal("parse_log_type"), NRet: 4, Protect: true}, tb); err != nil {
		return
	}
	logTypeRet, _sourceRet, _sourceJsonRet, errString := L.Get(-4), L.Get(-3), L.Get(-2), L.Get(-1)
	L.Pop(3)
	if errString.String() != "" {
		err = errors.New(errString.String())
		return
	}
	if res, ok := logTypeRet.(lua.LString); !ok {
		err = errors.New("logTypeRet.(lua.LString) -> !ok")
		return
	} else {
		logType = service.LogType(res)
	}
	if res, ok := _sourceRet.(*lua.LTable); !ok {
		err = errors.New("_sourceRet.(*lua.LTable) -> !ok")
		return
	} else {
		_sourceTable = res
	}
	if res, ok := _sourceJsonRet.(lua.LString); !ok {
		err = errors.New("_sourceJsonRet.(lua.LString) -> !ok")
		return
	} else {
		_sourceString = string(res)
	}
	return
}

func parseLog[T service.Log](parser *LogParser, parserFields []config.ParserField, log []byte, source map[string]interface{}, _sourceTable *lua.LTable, res T) (err error) {
	g := gjson.ParseBytes(log)
	targetJSON := ""
	for _, field := range parserFields {
		handleParserField(parser, &field, g, &targetJSON, _sourceTable)
	}
	if err = json.Unmarshal(unsafe.Slice(unsafe.StringData(targetJSON), len(targetJSON)), &res); err != nil {
		return
	}
	res.SetSource(source)
	res.Finish()
	return
}

func handleParserField(parser *LogParser, field *config.ParserField, g gjson.Result, targetJSON *string, _sourceTable *lua.LTable) {
	switch field.Type {
	case constant.ParserFieldTypeReplacements:
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
	case constant.ParserFieldTypeLua:
		L, fn := parser.LuaState.GetLuaState()
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
			var newValue interface{}
			switch field.LuaReturn {
			case constant.ParserFieldReturnString:
				if res, ok := ret.(lua.LString); !ok {
					continue
				} else {
					newValue = string(res)
				}
			case constant.ParserFieldReturnNumber:
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
