package logic_log_parser

import (
	"encoding/json"
	"sort"
	"unsafe"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
	"github.com/tidwall/gjson"
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
			var logInterface interface{}
			data, _ := json.Marshal(hit)

			log := new(service.AccessLog)
			_ = parseLog(backend.ParserFields.AccessLog, data, log)

			logInterface = log
			logs = append(logs, service.LogItem{
				Timestamp: gotil.ParseTime(log.Time),
				Storage:   backend.Name,
				LogType:   service.LogTypeAccessLog,
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

func parseLog[T service.Log](parserFields []config.ParserField, log []byte, res T) (err error) {
	g := gjson.ParseBytes(log)
	targetJSON := ""
	for _, field := range parserFields {
		field.Handle(g, &targetJSON, string(log))
	}
	if err = json.Unmarshal(unsafe.Slice(unsafe.StringData(targetJSON), len(targetJSON)), &res); err != nil {
		return
	}
	source := map[string]interface{}{}
	if err = json.Unmarshal(log, &source); err != nil {
		return
	}
	res.SetSource(source)
	return
}
