package logic

import (
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

type QueryBuilder interface {
	SearchQueryElastic(backend config.Backend, req service.SearchRequest) (
		queries map[string]elastic.Query,
		aggregations map[string]elastic.Aggregation,
		err error,
	)
	SearchQuerySLS(backend config.Backend, req service.SearchRequest) (
		queries map[string]string,
		err error,
	)
}
