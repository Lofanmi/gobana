package logic

import (
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

type LogParser interface {
	ParseElastic(backend config.Backend, m map[string]*elastic.SearchResult) (total int, logs service.LogItems, err error)
}
