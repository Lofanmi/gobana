package logic_log_parser

import (
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

var (
	_ logic.LogParser = &LogParser{}
)

// LogParser
// @autowire(logic.LogParser,set=logics)
type LogParser struct{}

func (s *LogParser) ParseElastic(backend config.Backend, m map[string]*elastic.SearchResult) (total int, logs []service.LogItem, err error) {
	total = s.parseElasticTotal(m)
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
