package svc_query_builder

import (
	"strings"

	"github.com/Lofanmi/gobana/service"
)

func (s *QueryBuilder) queryByHumanSLS(indexName string, query service.QueryByHuman) (queryRes string) {
	if len(query.Or) <= 0 && len(query.Must) <= 0 && len(query.MustNot) <= 0 {
		return
	}
	index, exist := s.Indexes[indexName]
	if !exist {
		return
	}
	var searchConditions, fuzzyConditions []string
	mainQuery := new(SlsQuery)
	searchConditions, fuzzyConditions = quoteConditions(query.Must)
	mainQuery.PrepareSearchConditions(searchConditions, operatorAnd, false)
	mainQuery.PrepareFuzzyConditions(index.DefaultFields, fuzzyConditions, operatorAnd, false)
	searchConditions, fuzzyConditions = quoteConditions(query.Or)
	mainQuery.PrepareSearchConditions(searchConditions, operatorOr, false)
	mainQuery.PrepareFuzzyConditions(index.DefaultFields, fuzzyConditions, operatorAnd, false)
	searchConditions, fuzzyConditions = quoteConditions(query.MustNot)
	mainQuery.PrepareSearchConditions(searchConditions, operatorAnd, true)
	mainQuery.PrepareFuzzyConditions(index.DefaultFields, fuzzyConditions, operatorAnd, true)
	if mainQuery.Empty() {
		return
	}
	index.BuildInQuery.Must.FieldConditions(func(field string, conditions []string) bool {
		searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
		mainQuery.PrepareSearchConditions(searchConditions2, operatorAnd, false)
		mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorAnd, false)
		return true
	})
	index.BuildInQuery.Or.FieldConditions(func(field string, conditions []string) bool {
		searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
		mainQuery.PrepareSearchConditions(searchConditions2, operatorOr, false)
		mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorOr, false)
		return true
	})
	index.BuildInQuery.MustNot.FieldConditions(func(field string, conditions []string) bool {
		searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
		mainQuery.PrepareSearchConditions(searchConditions2, operatorAnd, true)
		mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorAnd, true)
		return true
	})
	queryRes = mainQuery.String()
	return
}

func (s *QueryBuilder) queryBySlsQuerySls(query service.QueryBySLSQuery) (queryRes string) {
	queryRes = strings.TrimSpace(query.SlsQuery)
	return
}

// ---------------------------------------------------------------------------------------------------------------------

func quoteConditions(conditions []string) (searchConditions, fuzzyConditions []string) {
	for _, condition := range conditions {
		if condition = strings.TrimSpace(condition); condition == "" {
			continue
		}
		condition = strings.ReplaceAll(condition, `'`, `\'`)
		i := strings.IndexByte(condition, '%')
		if i == 0 || i == len(condition)-1 {
			fuzzyConditions = append(fuzzyConditions, `"`+condition+`"`)
		} else {
			searchConditions = append(searchConditions, `"`+condition+`"`)
		}
	}
	return
}
