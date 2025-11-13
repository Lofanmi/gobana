package svc_query_builder

import (
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

var (
	_ service.QueryBuilder = &QueryBuilder{}
)

// QueryBuilder
// @autowire(service.QueryBuilder,set=service)
type QueryBuilder struct {
	ApplicationConfig config.Application
	Indexes           config.Indexes
}

func (s *QueryBuilder) BuildSls(indexName string, req service.SearchRequest) (queryRes string, err error) {
	switch req.QueryBy {
	case service.QueryTypeByHuman:
		var q service.QueryByHuman
		if err = gotil.JsonAs(&req.Query, &q); err != nil {
			return
		}
		queryRes = s.queryByHumanSLS(indexName, q)
	case service.QueryTypeBySLSQuery:
		var q service.QueryBySLSQuery
		if err = gotil.JsonAs(&req.Query, &q); err != nil {
			return
		}
		queryRes = s.queryBySlsQuerySls(q)
	}
	return
}

func (s *QueryBuilder) BuildElastic(indexName string, req service.SearchRequest) (queryRes elastic.Query, aggregationRes elastic.Aggregation, err error) {
	switch req.QueryBy {
	case service.QueryTypeByHuman:
		var q service.QueryByHuman
		if err = gotil.JsonAs(&req.Query, &q); err != nil {
			return
		}
		queryRes, aggregationRes = s.queryByHumanElastic(indexName, req, q)
	case service.QueryTypeByLucene:
		var q service.QueryByLucene
		if err = gotil.JsonAs(&req.Query, &q); err != nil {
			return
		}
		queryRes, aggregationRes = s.queryByLuceneElastic(indexName, req, q)
	}
	return
}
