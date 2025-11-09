package service

import (
	"github.com/olivere/elastic/v7"
)

type AggregationParser interface {
	ParseElastic(timeA, timeB, interval int64, m map[string]*elastic.SearchResult) (xAxis []string, yAxis []int64, err error)
	ParseSLS(timeA, timeB, interval int64, m map[string]SlsSearchResult) (xAxis []string, yAxis []int64, err error)
}
