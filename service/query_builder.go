package service

import (
	"github.com/olivere/elastic/v7"
)

type QueryBuilder interface {
	BuildSls(backendName string, req SearchRequest) (queries map[string]string, err error)
	BuildElastic(backendName string, req SearchRequest) (queries map[string]elastic.Query, aggregations map[string]elastic.Aggregation, err error)
}
