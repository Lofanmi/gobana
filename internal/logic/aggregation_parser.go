package logic

import (
	"github.com/olivere/elastic/v7"
)

type AggregationParser interface {
	ParseElastic(timeA, timeB, interval int64, m map[string]*elastic.SearchResult) (xAxis []string, yAxis []int64, err error)
}
