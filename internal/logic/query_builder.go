package logic

import (
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

type QueryBuilder interface {
	SearchQueryElastic(backend config.Backend, req service.SearchRequest) (query map[string]elastic.Query, trackTotalHits bool, err error)
}
