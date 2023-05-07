package logic

import (
	"time"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

type Parser interface {
	ParseLogType() (logType service.LogType)
	ParseTime() (t time.Time)
	ParseHostname() (hostname string)
}

type LogParser interface {
	ParseElastic(backend config.Backend, m map[string]*elastic.SearchResult) (total int, logs []service.LogItem, err error)
}
