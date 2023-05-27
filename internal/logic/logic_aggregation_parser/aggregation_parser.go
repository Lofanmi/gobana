package logic_aggregation_parser

import (
	"strconv"

	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"github.com/wangjia184/sortedset"
)

var (
	_ logic.AggregationParser = &AggregationParser{}
)

// AggregationParser
// @autowire(logic.AggregationParser,set=logics)
type AggregationParser struct {
}

func (s *AggregationParser) ParseElastic(timeA, timeB, interval int64, m map[string]*elastic.SearchResult) (xAxis []string, yAxis []int64, err error) {
	timeRange := (timeB - timeA) / 1000
	points := timeRange / interval
	if timeRange%interval != 0 {
		points += 1
	}
	set := sortedset.New()
	for _, result := range m {
		if dateHistogram, ok := result.Aggregations.DateHistogram("charts"); ok {
			for _, bucket := range dateHistogram.Buckets {
				score := int(bucket.Key)
				keyName := strconv.Itoa(score)
				value := set.GetByKey(keyName)
				sum := bucket.DocCount
				if value != nil {
					sum += cast.ToInt64(value.Value)
				}
				set.AddOrUpdate(keyName, sortedset.SCORE(score), sum)
			}
		}
	}
	intervalSecond := interval * 1000
	xAxis = make([]string, 0, points)
	yAxis = make([]int64, 0, points)
	for i := int64(0); i < points; i++ {
		a := timeA + i*intervalSecond
		b := a + intervalSecond
		n := set.GetByScoreRange(sortedset.SCORE(a), sortedset.SCORE(b), &sortedset.GetByScoreRangeOptions{ExcludeEnd: true})
		value := int64(0)
		for _, node := range n {
			value += cast.ToInt64(node.Value)
		}
		xAxis = append(xAxis, gotil.FormatMilliSecond(a))
		yAxis = append(yAxis, value)
	}
	return
}
