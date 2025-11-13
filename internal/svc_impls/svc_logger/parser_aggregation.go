package svc_logger

import (
	"strconv"

	"github.com/Lofanmi/cmap"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/service"
	"github.com/spf13/cast"
	"github.com/wangjia184/sortedset"
)

func (s *Service) ParseAggregation(timeA, timeB, interval int64, m *cmap.Map[string, service.SearchResult]) (xAxis []string, yAxis []int64, err error) {
	indexes := m.Keys()
	xAxis, yAxis, err = s.parseAggregation(timeA, timeB, interval, func(set *sortedset.SortedSet) {
		for _, index := range indexes {
			searchResult, exist := m.Get(index)
			if !exist {
				continue
			}
			switch searchResult.ProviderType {
			case service.ProviderTypeSls:
				ms := timeA % 1000
				s.parseAggregationSls(searchResult.SearchResultSls, set, ms)
			case service.ProviderTypeElasticsearch:
				s.parseAggregationElastic(searchResult.SearchResultElastic, set)
			}
		}
	})
	return
}

func (s *Service) parseAggregationSls(sr service.SearchResultSls, set *sortedset.SortedSet, ms int64) {
	result := sr.ResponseAggregation
	if result == nil || len(result.Histograms) <= 0 {
		return
	}
	for _, singleHistogram := range result.Histograms {
		score := int(singleHistogram.From*1000 + ms)
		keyName := strconv.Itoa(score)
		value := set.GetByKey(keyName)
		sum := singleHistogram.Count
		if value != nil {
			sum += cast.ToInt64(value.Value)
		}
		set.AddOrUpdate(keyName, sortedset.SCORE(score), sum)
	}
	return
}

func (s *Service) parseAggregationElastic(sr service.SearchResultElastic, set *sortedset.SortedSet) {
	dateHistogram, ok := sr.Response.Aggregations.DateHistogram("charts")
	if !ok {
		return
	}
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
	return
}

func (s *Service) parseAggregation(timeA, timeB, interval int64, fn func(set *sortedset.SortedSet)) (xAxis []string, yAxis []int64, err error) {
	timeRange := (timeB - timeA) / 1000
	points := timeRange / interval
	if timeRange%interval != 0 {
		points += 1
	}
	set := sortedset.New()
	fn(set)
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
