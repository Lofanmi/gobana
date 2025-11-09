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
	Backends          config.Backends
	Indexes           config.Indexes
}

func (s *QueryBuilder) BuildSls(backendName string, req service.SearchRequest) (queries map[string]string, err error) {
	switch req.QueryBy {
	case service.QueryTypeByHuman:
		var q service.QueryByHuman
		if err = gotil.JsonAs(&req.Query, &q); err != nil {
			return
		}
		queries = s.queryByHumanSLS(s.Backends[backendName], req, q)
	case service.QueryTypeBySLSQuery:
		var q service.QueryBySLSQuery
		if err = gotil.JsonAs(&req.Query, &q); err != nil {
			return
		}
		queries = s.queryBySLSQuerySLS(s.Backends[backendName], req, q)
	}
	return
}

func (s *QueryBuilder) BuildElastic(backendName string, req service.SearchRequest) (queries map[string]elastic.Query, aggregations map[string]elastic.Aggregation, err error) {
	switch req.QueryBy {
	case service.QueryTypeByHuman:
		var q service.QueryByHuman
		if err = gotil.JsonAs(&req.Query, &q); err != nil {
			return
		}
		queries, aggregations = s.queryByHumanElastic(s.Backends[backendName], req, q)
	case service.QueryTypeByLucene:
		var q service.QueryByLucene
		if err = gotil.JsonAs(&req.Query, &q); err != nil {
			return
		}
		queries, aggregations = s.queryByLuceneElastic(s.Backends[backendName], req, q)
	}
	return
}
