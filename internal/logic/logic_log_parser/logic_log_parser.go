package logic_log_parser

import (
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

type LogParser struct{}

func (s *LogParser) ParseElastic(backend config.Backend, m map[string]*elastic.SearchResult) (total int, logs []service.LogItem, err error) {
	// TODO implement me
	panic("implement me")
}
