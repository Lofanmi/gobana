package svc_query_builder

import (
	"strings"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/service"
)

func (s *QueryBuilder) queryByHumanSLS(backendConfig config.BackendConfig, req service.SearchRequest, query service.QueryByHuman) (queries map[string]string) {
	if len(query.Or) <= 0 && len(query.Must) <= 0 && len(query.MustNot) <= 0 {
		return
	}

	var searchConditions, fuzzyConditions []string
	queries = map[string]string{}
	for _, indexName := range backendConfig.Indexes {
		defaultFields := gotil.OrSliceDefault(backendConfig.DefaultFields[indexName], backendConfig.DefaultFields[service.DefaultValue])
		buildInQuery, ok := backendConfig.BuildInQuery[indexName]
		if !ok {
			buildInQuery = backendConfig.BuildInQuery[service.DefaultValue]
		}

		mainQuery := new(SlsQuery)
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
		buildInQuery.Must.FieldConditions(func(field string, conditions []string) bool {
			searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
			mainQuery.PrepareSearchConditions(searchConditions2, operatorAnd, false)
			mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorAnd, false)
			return true
		})
		buildInQuery.Or.FieldConditions(func(field string, conditions []string) bool {
			searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
			mainQuery.PrepareSearchConditions(searchConditions2, operatorOr, false)
			mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorOr, false)
			return true
		})
		buildInQuery.MustNot.FieldConditions(func(field string, conditions []string) bool {
			searchConditions2, fuzzyConditions2 := quoteConditions(conditions)
			mainQuery.PrepareSearchConditions(searchConditions2, operatorAnd, true)
			mainQuery.PrepareFuzzyConditions([]string{field}, fuzzyConditions2, operatorAnd, true)
			return true
		})
		queries[indexName] = mainQuery.String()
	}

	return
}

func (s *QueryBuilder) queryBySLSQuerySLS(backendConfig config.BackendConfig, req service.SearchRequest, query service.QueryBySLSQuery) (queries map[string]string) {
	query.SlsQuery = strings.TrimSpace(query.SlsQuery)
	if len(query.SlsQuery) <= 0 {
		return
	}

	queries = map[string]string{}
	for _, indexName := range backendConfig.Indexes {
		queries[indexName] = query.SlsQuery
	}

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
