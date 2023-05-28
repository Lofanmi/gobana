package logic_query_builder

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/constant"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
)

var (
	_ logic.QueryBuilder = &QueryBuilder{}
)

// QueryBuilder
// @autowire(logic.QueryBuilder,set=logics)
type QueryBuilder struct{}

func (s *QueryBuilder) SearchQueryElastic(backend config.Backend, req service.SearchRequest) (
	queries map[string]elastic.Query,
	aggregations map[string]elastic.Aggregation,
	err error,
) {
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
		queries, aggregations = s.queryByHumanElastic(backend, req, q)
	case service.QueryTypeByLucene:
		var q service.QueryByLucene
		if err = json.Unmarshal(data, &q); err != nil {
			return
		}
		queries, aggregations = s.queryByLuceneElastic(backend, req, q)
	}
	return
}

func (s *QueryBuilder) queryByHumanElastic(backend config.Backend, req service.SearchRequest, query service.QueryByHuman) (
	queries map[string]elastic.Query,
	aggregations map[string]elastic.Aggregation,
) {
	if len(query.Or) <= 0 && len(query.Must) <= 0 && len(query.MustNot) <= 0 {
		return
	}
	queries = map[string]elastic.Query{}
	aggregations = map[string]elastic.Aggregation{}
	indexList := backend.MultiSearch[req.Storage].IndexList
	for _, index := range indexList {
		defaultFields, ok := backend.DefaultFields[index]
		if !ok {
			defaultFields = backend.DefaultFields[constant.DefaultValue]
		}
		timeField := backend.TimeField[index]
		if timeField == "" {
			timeField = backend.TimeField[constant.DefaultValue]
		}
		if timeField == "" {
			timeField = constant.AtTimestamp
		}
		timezone := backend.Timezone[index]
		if timezone == "" {
			timezone = backend.Timezone[constant.DefaultValue]
		}
		esMainQuery := elastic.NewBoolQuery()
		emptyCondition := true
		TimeQuery(timeField, req.TimeA, req.TimeB, func(query elastic.Query) { esMainQuery.Must(query) })
		OrQueries(defaultFields, query.Or, &emptyCondition, func(orQueries []elastic.Query) {
			esMainQuery.Should(orQueries...).MinimumNumberShouldMatch(1)
		})
		MustOrMustNotQueries(defaultFields, query.Must, &emptyCondition, func(query elastic.Query) { esMainQuery.Must(query) })
		MustOrMustNotQueries(defaultFields, query.MustNot, &emptyCondition, func(query elastic.Query) { esMainQuery.MustNot(query) })
		if emptyCondition {
			queries[index] = esMainQuery
			continue
		}
		buildInQueries, ok := backend.BuildInQueries[index]
		if !ok {
			buildInQueries = backend.BuildInQueries[constant.DefaultValue]
		}
		MustOrMustNotBuildInQueryEntry(buildInQueries.Must, func(query elastic.Query) { esMainQuery.Must(query) })
		MustOrMustNotBuildInQueryEntry(buildInQueries.MustNot, func(query elastic.Query) { esMainQuery.MustNot(query) })
		OrBuildInQueryEntry(buildInQueries.Or, func(orQueries []elastic.Query) {
			esMainQuery.Should(orQueries...).MinimumNumberShouldMatch(1)
		})
		queries[index] = esMainQuery
		if req.ChartVisible {
			aggregations[index] = elastic.NewDateHistogramAggregation().
				Field(timeField).
				FixedInterval(strconv.Itoa(int(req.ChartInterval)) + "s").
				TimeZone(timezone).
				MinDocCount(0)
		}
	}
	return
}

func (s *QueryBuilder) queryByLuceneElastic(backend config.Backend, req service.SearchRequest, query service.QueryByLucene) (
	queries map[string]elastic.Query,
	aggregations map[string]elastic.Aggregation,
) {
	if len(query.Lucene) <= 0 {
		return
	}
	queries = map[string]elastic.Query{}
	aggregations = map[string]elastic.Aggregation{}
	indexList := backend.MultiSearch[req.Storage].IndexList
	for _, index := range indexList {
		timeField := backend.TimeField[index]
		if timeField == "" {
			timeField = backend.TimeField[constant.DefaultValue]
		}
		if timeField == "" {
			timeField = constant.AtTimestamp
		}
		timezone := backend.Timezone[index]
		if timezone == "" {
			timezone = backend.Timezone[constant.DefaultValue]
		}
		esMainQuery := elastic.NewBoolQuery()
		TimeQuery(timeField, req.TimeA, req.TimeB, func(query elastic.Query) { esMainQuery.Filter(query) })
		esMainQuery.Filter(elastic.NewQueryStringQuery(query.Lucene))
		queries[index] = esMainQuery
		if req.ChartVisible {
			aggregations[index] = elastic.NewDateHistogramAggregation().
				Field(timeField).
				FixedInterval(strconv.Itoa(int(req.ChartInterval)) + "s").
				TimeZone(timezone).
				MinDocCount(0)
		}
	}
	return
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

func (s *QueryBuilder) SearchQuerySLS(backend config.Backend, req service.SearchRequest) (
	queries map[string]string,
	err error,
) {
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
		queries = s.queryByHumanSLS(backend, req, q)
	case service.QueryTypeBySLSQuery:
		var q service.QueryBySLSQuery
		if err = json.Unmarshal(data, &q); err != nil {
			return
		}
		queries = s.queryBySLSQuerySLS(backend, req, q)
	}
	return
}

func (s *QueryBuilder) queryByHumanSLS(backend config.Backend, req service.SearchRequest, query service.QueryByHuman) (
	queries map[string]string,
) {
	if len(query.Or) <= 0 && len(query.Must) <= 0 && len(query.MustNot) <= 0 {
		return
	}
	queries = map[string]string{}
	indexList := backend.MultiSearch[req.Storage].IndexList
	for _, index := range indexList {
		defaultFields, ok := backend.DefaultFields[index]
		if !ok {
			defaultFields = backend.DefaultFields[constant.DefaultValue]
		}
		var merge []string
		if q := AndOrNotQueries(defaultFields, query.Must, " AND ", false); q != "" {
			merge = append(merge, q)
		}
		if q := AndOrNotQueries(defaultFields, query.Or, " OR ", false); q != "" {
			merge = append(merge, q)
		}
		if q := AndOrNotQueries(defaultFields, query.MustNot, " AND ", true); q != "" {
			merge = append(merge, q)
		}
		if len(merge) <= 0 {
			continue
		}
		buildInQueries, ok := backend.BuildInQueries[index]
		if !ok {
			buildInQueries = backend.BuildInQueries[constant.DefaultValue]
		}
		if q := AndOrNotBuildInQueryEntry(buildInQueries.Must, " AND ", false); q != "" {
			merge = append(merge, q)
		}
		if q := AndOrNotBuildInQueryEntry(buildInQueries.Or, " OR ", false); q != "" {
			merge = append(merge, q)
		}
		if q := AndOrNotBuildInQueryEntry(buildInQueries.MustNot, " AND ", true); q != "" {
			merge = append(merge, q)
		}
		queries[index] = strings.Join(merge, " AND ")
	}
	return
}

func AndOrNotQueries(fields, conditions []string, op string, not bool) string {
	var quoteConditions []string
	for _, condition := range conditions {
		if condition != "" {
			condition = strings.ReplaceAll(condition, `'`, `\'`)
			quoteConditions = append(quoteConditions, condition)
		}
	}
	if len(conditions) <= 0 {
		return ""
	}
	var res []string
	for _, field := range fields {
		operator := op
		if field != "__raw__" {
			operator = " OR "
		}
		var subQueries []string
		for _, condition := range quoteConditions {
			if not {
				subQueries = append(subQueries, fmt.Sprintf(`("%s" not like '%%%s%%')`, field, condition))
			} else {
				subQueries = append(subQueries, fmt.Sprintf(`"%s" like '%%%s%%'`, field, condition))
			}
		}
		if len(subQueries) <= 0 {
			continue
		}
		res = append(res, "("+strings.Join(subQueries, operator)+")")
	}
	if len(res) <= 0 {
		return ""
	}
	return "(" + strings.Join(res, " OR ") + ")"
}

func AndOrNotBuildInQueryEntry(items []config.BuildInQueryEntry, op string, not bool) string {
	var subQueries []string
	for _, item := range items {
		if !item.Always {
			continue
		}
		if q := AndOrNotQueries([]string{item.Field}, cast.ToStringSlice(item.Values), op, not); q != "" {
			subQueries = append(subQueries, q)
		}
	}
	if len(subQueries) <= 0 {
		return ""
	}
	return "(" + strings.Join(subQueries, " OR ") + ")"
}

func (s *QueryBuilder) queryBySLSQuerySLS(backend config.Backend, req service.SearchRequest, query service.QueryBySLSQuery) (
	queries map[string]string,
) {
	if len(query.SQL) <= 0 {
		return
	}
	queries = map[string]string{}
	indexList := backend.MultiSearch[req.Storage].IndexList
	for _, index := range indexList {
		queries[index] = query.SQL
	}
	return
}
