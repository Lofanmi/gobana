package logic_log_parser

import (
	"encoding/json"
	"errors"
	"sort"
	"unsafe"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
	"github.com/tidwall/gjson"
	lua "github.com/yuin/gopher-lua"
)

var (
	_ logic.LogParser = &LogParser{}
)

// LogParser
// @autowire(logic.LogParser,set=logics)
type LogParser struct{}

func (s *LogParser) ParseElastic(backend config.Backend, m map[string]*elastic.SearchResult) (total int, logs service.LogItems, err error) {
	if total = parseElasticTotal(m); total <= 0 {
		return
	}
	logs = make([]service.LogItem, 0, total)
	for _, result := range m {
		if result == nil || result.Hits == nil || result.Hits.TotalHits == nil {
			continue
		}
		for _, hit := range result.Hits.Hits {
			data, _ := json.Marshal(hit)
			source := map[string]interface{}{}
			if err = json.Unmarshal(data, &source); err != nil {
				return
			}
			tb := gotil.MapToTable(source)
			var logInterface interface{}
			logTime := ""
			logType, e := parseLogType(backend, tb)
			switch logType {
			case service.LogTypeAccessLog:
				log := new(service.AccessLog)
				_ = parseLog(backend.ParserFields.AccessLog, data, source, tb, log)
				logInterface, logTime = log, log.Time
			case service.LogTypeJsonLog:
				log := new(service.JsonLog)
				_ = parseLog(backend.ParserFields.JsonLog, data, source, tb, log)
				logInterface, logTime = log, log.Time
			case service.LogTypeStringLog:
				log := new(service.StringLog)
				_ = parseLog(backend.ParserFields.StringLog, data, source, tb, log)
				logInterface, logTime = log, log.Time
			default:
				_ = e
				continue
			}
			logs = append(logs, service.LogItem{
				Timestamp: gotil.ParseTime(logTime),
				Storage:   hit.Index,
				LogType:   logType,
				Log:       logInterface,
			})
		}
	}
	sort.Sort(logs)
	return
}

func parseElasticTotal(m map[string]*elastic.SearchResult) (total int) {
	for _, result := range m {
		if result == nil || result.Hits == nil || result.Hits.TotalHits == nil {
			continue
		}
		total += int(result.Hits.TotalHits.Value)
	}
	return
}

func parseLogType(backend config.Backend, tb *lua.LTable) (logType service.LogType, err error) {
	L, fn := config.GetLuaState()
	defer fn()
	if err = L.DoString(backend.ParserLogType); err != nil {
		return
	}
	if err = L.CallByParam(lua.P{Fn: L.GetGlobal("parse_log_type"), NRet: 2, Protect: true}, tb); err != nil {
		return
	}
	ret, errString := L.Get(-2), L.Get(-1)
	L.Pop(2)
	if errString.String() != "" {
		return
	}
	if res, ok := ret.(lua.LString); !ok {
		err = errors.New("ret.(lua.LString) -> !ok")
	} else {
		logType = service.LogType(res)
	}
	return
}

func parseLog[T service.Log](parserFields []config.ParserField, log []byte, source map[string]interface{}, tb *lua.LTable, res T) (err error) {
	g := gjson.ParseBytes(log)
	targetJSON := ""
	for _, field := range parserFields {
		field.Handle(g, &targetJSON, tb)
	}
	if err = json.Unmarshal(unsafe.Slice(unsafe.StringData(targetJSON), len(targetJSON)), &res); err != nil {
		return
	}
	res.SetSource(source)
	return
}
