package svc_query_builder

import (
	"strconv"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

func (s *QueryBuilder) queryByHumanElastic(indexName string, req service.SearchRequest, query service.QueryByHuman) (queryRes elastic.Query, aggregationRes elastic.Aggregation) {
	if len(query.Or) <= 0 && len(query.Must) <= 0 && len(query.MustNot) <= 0 {
		return
	}
	index, exist := s.Indexes[indexName]
	if !exist {
		return
	}
	indexMeta := index.Meta.IndexMetaForElastic

	esMainQuery := elastic.NewBoolQuery()
	emptyCondition := true
	TimeQuery(indexMeta.TimeField, req.TimeA, req.TimeB, func(query elastic.Query) { esMainQuery.Must(query) })
	OrQueries(index.DefaultFields, query.Or, &emptyCondition, func(orQueries []elastic.Query) {
		esMainQuery.Should(orQueries...).MinimumNumberShouldMatch(1)
	})
	MustOrMustNotQueries(index.DefaultFields, query.Must, &emptyCondition, func(query elastic.Query) { esMainQuery.Must(query) })
	MustOrMustNotQueries(index.DefaultFields, query.MustNot, &emptyCondition, func(query elastic.Query) { esMainQuery.MustNot(query) })
	if emptyCondition {
		queryRes = esMainQuery
		return
	}
	MustOrMustNotBuildInQueryEntry(index.BuildInQuery.Must, func(query elastic.Query) { esMainQuery.Must(query) })
	MustOrMustNotBuildInQueryEntry(index.BuildInQuery.MustNot, func(query elastic.Query) { esMainQuery.MustNot(query) })
	OrBuildInQueryEntry(index.BuildInQuery.Or, func(orQueries []elastic.Query) {
		esMainQuery.Should(orQueries...).MinimumNumberShouldMatch(1)
	})

	queryRes = esMainQuery
	if req.ChartVisible {
		aggregationRes = elastic.NewDateHistogramAggregation().Field(indexMeta.TimeField).FixedInterval(strconv.Itoa(req.ChartInterval) + "s").TimeZone(indexMeta.Timezone).MinDocCount(0)
	}
	return
}

func (s *QueryBuilder) queryByLuceneElastic(indexName string, req service.SearchRequest, query service.QueryByLucene) (queryRes elastic.Query, aggregationRes elastic.Aggregation) {
	if len(query.Lucene) <= 0 {
		return
	}
	index, exist := s.Indexes[indexName]
	if !exist {
		return
	}
	indexMeta := index.Meta.IndexMetaForElastic

	esMainQuery := elastic.NewBoolQuery()
	TimeQuery(indexMeta.TimeField, req.TimeA, req.TimeB, func(query elastic.Query) { esMainQuery.Filter(query) })
	esMainQuery.Filter(elastic.NewQueryStringQuery(query.Lucene))

	queryRes = esMainQuery
	if req.ChartVisible {
		aggregationRes = elastic.NewDateHistogramAggregation().Field(indexMeta.TimeField).FixedInterval(strconv.Itoa(req.ChartInterval) + "s").TimeZone(indexMeta.Timezone).MinDocCount(0)
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------

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

func TimeQuery(timeField string, timeA, timeB int64, fn func(query elastic.Query)) {
	query := elastic.NewRangeQuery(timeField).
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
