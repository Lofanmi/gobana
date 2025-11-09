package svc_query_builder

import (
	"strconv"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
)

func (s *QueryBuilder) queryByHumanElastic(backendConfig config.BackendConfig, req service.SearchRequest, query service.QueryByHuman) (queries map[string]elastic.Query, aggregations map[string]elastic.Aggregation) {
	if len(query.Or) <= 0 && len(query.Must) <= 0 && len(query.MustNot) <= 0 {
		return
	}

	queries = map[string]elastic.Query{}
	aggregations = map[string]elastic.Aggregation{}

	for _, indexName := range backendConfig.Indexes {
		defaultFields := gotil.OrSliceDefault(backendConfig.DefaultFields[indexName], backendConfig.DefaultFields[service.DefaultValue])
		timeField := gotil.OrDefault(backendConfig.TimeField[indexName], backendConfig.TimeField[service.DefaultValue], service.AtTimestamp)
		timezone := gotil.OrDefault(backendConfig.Timezone[indexName], backendConfig.Timezone[service.DefaultValue], s.ApplicationConfig.Timezone)

		esMainQuery := elastic.NewBoolQuery()
		emptyCondition := true
		TimeQuery(timeField, req.TimeA, req.TimeB, func(query elastic.Query) { esMainQuery.Must(query) })
		OrQueries(defaultFields, query.Or, &emptyCondition, func(orQueries []elastic.Query) {
			esMainQuery.Should(orQueries...).MinimumNumberShouldMatch(1)
		})
		MustOrMustNotQueries(defaultFields, query.Must, &emptyCondition, func(query elastic.Query) { esMainQuery.Must(query) })
		MustOrMustNotQueries(defaultFields, query.MustNot, &emptyCondition, func(query elastic.Query) { esMainQuery.MustNot(query) })
		if emptyCondition {
			queries[indexName] = esMainQuery
			continue
		}

		buildInQuery, ok := backendConfig.BuildInQuery[indexName]
		if !ok {
			buildInQuery = backendConfig.BuildInQuery[service.DefaultValue]
		}
		MustOrMustNotBuildInQueryEntry(buildInQuery.Must, func(query elastic.Query) { esMainQuery.Must(query) })
		MustOrMustNotBuildInQueryEntry(buildInQuery.MustNot, func(query elastic.Query) { esMainQuery.MustNot(query) })
		OrBuildInQueryEntry(buildInQuery.Or, func(orQueries []elastic.Query) {
			esMainQuery.Should(orQueries...).MinimumNumberShouldMatch(1)
		})
		queries[indexName] = esMainQuery

		if req.ChartVisible {
			aggregations[indexName] = elastic.NewDateHistogramAggregation().
				Field(timeField).
				FixedInterval(strconv.Itoa(req.ChartInterval) + "s").
				TimeZone(timezone).
				MinDocCount(0)
		}
	}

	return
}

func (s *QueryBuilder) queryByLuceneElastic(backendConfig config.BackendConfig, req service.SearchRequest, query service.QueryByLucene) (queries map[string]elastic.Query, aggregations map[string]elastic.Aggregation) {
	if len(query.Lucene) <= 0 {
		return
	}

	queries = map[string]elastic.Query{}
	aggregations = map[string]elastic.Aggregation{}

	for _, indexName := range backendConfig.Indexes {
		timeField := gotil.OrDefault(backendConfig.TimeField[indexName], backendConfig.TimeField[service.DefaultValue], service.AtTimestamp)
		timezone := gotil.OrDefault(backendConfig.Timezone[indexName], backendConfig.Timezone[service.DefaultValue], s.ApplicationConfig.Timezone)

		esMainQuery := elastic.NewBoolQuery()
		TimeQuery(timeField, req.TimeA, req.TimeB, func(query elastic.Query) { esMainQuery.Filter(query) })
		esMainQuery.Filter(elastic.NewQueryStringQuery(query.Lucene))
		queries[indexName] = esMainQuery

		if req.ChartVisible {
			aggregations[indexName] = elastic.NewDateHistogramAggregation().
				Field(timeField).
				FixedInterval(strconv.Itoa(req.ChartInterval) + "s").
				TimeZone(timezone).
				MinDocCount(0)
		}
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
