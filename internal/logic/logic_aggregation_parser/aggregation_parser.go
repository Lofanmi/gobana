package logic_aggregation_parser

import (
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/olivere/elastic/v7"
)

var (
	_ logic.AggregationParser = &AggregationParser{}
)

// AggregationParser
// @autowire(logic.AggregationParser,set=logics)
type AggregationParser struct {
}

func (s *AggregationParser) ParseElastic(timeA, timeB, interval int64, m map[string]*elastic.SearchResult) (xAxis []string, yAxis []int64, err error) {
	timeRange := timeB - timeA
	points := timeRange / interval
	if timeRange%interval != 0 {
		points += 1
	}
	return
}
