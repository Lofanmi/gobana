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

	var searchConditions, fuzzyConditions []string
	for _, index := range indexList {
		defaultFields, ok := backend.DefaultFields[index]
		if !ok {
			defaultFields = backend.DefaultFields[constant.DefaultValue]
		}
		buildInQueries, ok := backend.BuildInQueries[index]
		if !ok {
			buildInQueries = backend.BuildInQueries[constant.DefaultValue]
		}
		mainQuery := new(slsQuery)
		searchConditions, fuzzyConditions = quoteConditions(query.Must)
		mainQuery.PrepareSearchConditions(searchConditions, operatorAnd, false)
		mainQuery.PrepareFuzzyConditions(defaultFields, fuzzyConditions, operatorAnd, false)
		searchConditions, fuzzyConditions = quoteConditions(query.Or)
		mainQuery.PrepareSearchConditions(searchConditions, operatorOr, false)
		mainQuery.PrepareFuzzyConditions(defaultFields, fuzzyConditions, operatorAnd, false)
		searchConditions, fuzzyConditions = quoteConditions(query.MustNot)
		mainQuery.PrepareSearchConditions(searchConditions, operatorAnd, true)
		mainQuery.PrepareFuzzyConditions(defaultFields, fuzzyConditions, operatorAnd, true)
		if mainQuery.Empty() {
			continue
		}
		buildInQueries.Must.FieldConditions(func(field string, conditions []string) bool {
			searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
			mainQuery.PrepareSearchConditions(searchConditions2, operatorAnd, false)
			mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorAnd, false)
			return true
		})
		buildInQueries.Or.FieldConditions(func(field string, conditions []string) bool {
			searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
			mainQuery.PrepareSearchConditions(searchConditions2, operatorOr, false)
			mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorOr, false)
			return true
		})
		buildInQueries.MustNot.FieldConditions(func(field string, conditions []string) bool {
			searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
			mainQuery.PrepareSearchConditions(searchConditions2, operatorAnd, true)
			mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorAnd, true)
			return true
		})
		queries[index] = mainQuery.String()
	}
	return
}

const (
	operatorAnd = " AND "
	operatorOr  = " OR "
	operatorNot = " AND NOT "
)

type slsQuery struct {
	searchConditions []string
	fuzzyConditions  []string
}

func (s *slsQuery) Empty() bool {
	return len(s.fuzzyConditions) <= 0 && len(s.searchConditions) <= 0
}

func (s *slsQuery) String() string {
	if s.Empty() {
		return ""
	}
	var search string
	if len(s.searchConditions) <= 0 {
		search = "*"
	} else {
		search = strings.Join(s.searchConditions, operatorAnd)
	}
	var where string
	if len(s.fuzzyConditions) > 0 {
		where = " | SELECT * FROM log WHERE " + strings.Join(s.fuzzyConditions, operatorAnd)
	}
	return search + where
}

func (s *slsQuery) PrepareSearchConditions(conditions []string, operator string, not bool) {
	if len(conditions) <= 0 {
		return
	}
	var condition string
	if not {
		condition = "NOT " + strings.Join(conditions, operator)
	} else {
		condition = strings.Join(conditions, operator)
	}
	s.searchConditions = append(s.searchConditions, "("+condition+")")
}

func (s *slsQuery) PrepareFuzzyConditions(fields []string, conditions []string, operator string, not bool) {
	if len(fields) <= 0 || len(conditions) <= 0 {
		return
	}
	var res []string
	for _, condition := range conditions {
		var subQueries []string
		for _, field := range fields {
			if not {
				subQueries = append(subQueries, fmt.Sprintf(`("%s" not like %s)`, field, condition))
			} else {
				subQueries = append(subQueries, fmt.Sprintf(`("%s" like %s)`, field, condition))
			}
		}
		if len(subQueries) <= 0 {
			continue
		}
		res = append(res, "("+strings.Join(subQueries, operatorOr)+")")
	}
	s.fuzzyConditions = append(s.fuzzyConditions, "("+strings.Join(res, operator)+")")
}

func quoteConditions(conditions []string) (searchConditions, fuzzyConditions []string) {
	for _, condition := range conditions {
		if condition = strings.TrimSpace(condition); condition == "" {
			continue
		}
		condition = strings.ReplaceAll(condition, `'`, `\'`)
		i := strings.IndexByte(condition, '%')
		if i == 0 || i == len(condition)-1 {
			fuzzyConditions = append(fuzzyConditions, "'"+condition+"'")
		} else {
			searchConditions = append(searchConditions, "'"+condition+"'")
		}
	}
	return
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
