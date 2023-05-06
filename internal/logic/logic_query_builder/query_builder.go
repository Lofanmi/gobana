package logic_query_builder

import (
	"encoding/json"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

const (
	atTimestamp = "@timestamp"
)

var (
	_ logic.QueryBuilder = &QueryBuilder{}
)

// QueryBuilder
// @autowire(logic.QueryBuilder,set=logics)
type QueryBuilder struct{}

func (s *QueryBuilder) SearchQueryElastic(backend config.Backend, req service.SearchRequest) (queries []elastic.Query, err error) {
	var data []byte
	data, err = json.Marshal(req.Query)
	if err != nil {
		return
	}
	switch req.QueryBy {
	case service.QueryTypeByHuman:
		var q service.QueryByHuman
		if err = json.Unmarshal(data, &q); err != nil {
			return
		}
		queries = s.queryByHuman(backend, req, q)
	case service.QueryTypeByQueryString:
		var q service.QueryByQueryString
		if err = json.Unmarshal(data, &q); err != nil {
			return
		}
		queries = s.queryByQueryString(backend, req, q)
	}
	return
}

func (s *QueryBuilder) queryByHuman(backend config.Backend, req service.SearchRequest, query service.QueryByHuman) (queries []elastic.Query) {
	indexList := backend.MultiSearch[req.Storage].IndexList
	for _, index := range indexList {
		defaultFields := backend.DefaultFields[index]
		esMainQuery := elastic.NewBoolQuery()
		emptyCondition := true
		TimeQuery(req.TimeA, req.TimeB, func(query elastic.Query) { esMainQuery.Must(query) })
		OrQueries(defaultFields, query.Or, &emptyCondition, func(orQueries []elastic.Query) {
			esMainQuery.Should(orQueries...).MinimumNumberShouldMatch(1)
		})
		MustOrMustNotQueries(defaultFields, query.Must, &emptyCondition, func(query elastic.Query) { esMainQuery.Must(query) })
		MustOrMustNotQueries(defaultFields, query.MustNot, &emptyCondition, func(query elastic.Query) { esMainQuery.MustNot(query) })
		if emptyCondition {
			queries = append(queries, esMainQuery)
			continue
		}
		buildInQueries := backend.BuildInQueries[index]
		MustOrMustNotBuildInQueryEntry(buildInQueries.Must, func(query elastic.Query) { esMainQuery.Must(query) })
		MustOrMustNotBuildInQueryEntry(buildInQueries.MustNot, func(query elastic.Query) { esMainQuery.MustNot(query) })
		OrBuildInQueryEntry(buildInQueries.Or, func(orQueries []elastic.Query) {
			esMainQuery.Should(orQueries...).MinimumNumberShouldMatch(1)
		})
		queries = append(queries, esMainQuery)
	}
	return
}

func (s *QueryBuilder) queryByQueryString(backend config.Backend, req service.SearchRequest, query service.QueryByQueryString) (queries []elastic.Query) {
	return nil
}

func MustOrMustNotBuildInQueryEntry(items []config.BuildInQueryEntry, fn func(query elastic.Query)) {
	for _, item := range items {
		if !item.Always {
			continue
		}
		for _, value := range item.Values {
			query := elastic.NewMatchQuery(item.Field, value).Operator(item.Operator)
			fn(query)
		}
	}
}

func OrBuildInQueryEntry(items []config.BuildInQueryEntry, fn func(orQueries []elastic.Query)) {
	var queries []elastic.Query
	for _, item := range items {
		if !item.Always {
			continue
		}
		for _, value := range item.Values {
			query := elastic.NewMatchQuery(item.Field, value).Operator(item.Operator)
			queries = append(queries, query)
		}
	}
	if len(queries) > 0 {
		fn(queries)
	}
}

func TimeQuery(timeA, timeB int64, fn func(query elastic.Query)) {
	query := elastic.NewRangeQuery(atTimestamp).
		Gte(timeA).
		Lte(timeB).
		Format("epoch_millis")
	fn(query)
}

func OrQueries(defaultFields, ors []string, emptySearchHit *bool, fn func(orQueries []elastic.Query)) {
	var orQueries []elastic.Query
	for _, or := range ors {
		if or == "" {
			continue
		}
		var orSubQueries []elastic.Query
		for _, field := range defaultFields {
			orSubQuery := elastic.NewMatchPhraseQuery(field, or)
			orSubQueries = append(orSubQueries, orSubQuery)
		}
		orQuery := elastic.NewBoolQuery().Should(orSubQueries...).MinimumNumberShouldMatch(1)
		orQueries = append(orQueries, orQuery)
		*emptySearchHit = false
	}
	if len(orQueries) > 0 {
		fn(orQueries)
	}
}

func MustOrMustNotQueries(defaultFields, conditions []string, emptySearchHit *bool, fn func(query elastic.Query)) {
	for _, condition := range conditions {
		if condition == "" {
			continue
		}
		var subQueries []elastic.Query
		for _, field := range defaultFields {
			conditionSubQuery := elastic.NewMatchPhraseQuery(field, condition)
			subQueries = append(subQueries, conditionSubQuery)
		}
		query := elastic.NewBoolQuery().Should(subQueries...).MinimumNumberShouldMatch(1)
		fn(query)
		*emptySearchHit = false
	}
}
