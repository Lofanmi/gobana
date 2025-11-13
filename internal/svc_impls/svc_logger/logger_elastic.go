package svc_logger

import (
	"context"

	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

func (s *Service) doSearchIndexElastic(ctx context.Context, cli *elastic.Client, req service.SearchRequest, indexName string, query elastic.Query, aggregation elastic.Aggregation) (result service.SearchResultElastic, err error) {
	result.RawQuery, result.ErrorResponse = query.Source()
	if result.ErrorResponse != nil {
		return
	}
	index := s.Indexes[indexName]
	search := cli.Search()
	search.Index(index.Meta.IndexMetaForElastic.Indexes...).TrackTotalHits(req.TrackTotalHits).Query(query).Pretty(false).Version(true)
	search.From((req.PageNo - 1) * req.PageSize).Size(req.PageSize)
	if aggregation != nil {
		search.Aggregation("charts", aggregation)
	}
	for _, sortField := range index.Meta.IndexMetaForElastic.SortFields {
		search.Sort(sortField.Field, sortField.Ascending)
	}
	result.Response, result.ErrorResponse = search.Do(ctx)
	return
}
