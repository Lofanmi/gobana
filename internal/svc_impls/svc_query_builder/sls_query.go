package svc_query_builder

import (
	"fmt"
	"strings"
)

const (
	operatorAnd = " AND "
	operatorOr  = " OR "
	operatorNot = " AND NOT "
)

type SlsQuery struct {
	searchConditions []string
	fuzzyConditions  []string
}

func (s *SlsQuery) Empty() bool {
	return len(s.searchConditions) <= 0 && len(s.fuzzyConditions) <= 0
}

func (s *SlsQuery) String() string {
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

func (s *SlsQuery) PrepareSearchConditions(conditions []string, operator string, not bool) {
	if len(conditions) <= 0 {
		return
	}
	var condition string
	if not {
		operator = " AND NOT "
		condition = "NOT " + strings.Join(conditions, operator)
	} else {
		condition = strings.Join(conditions, operator)
	}
	s.searchConditions = append(s.searchConditions, "("+condition+")")
}

func (s *SlsQuery) PrepareFuzzyConditions(fields []string, conditions []string, operator string, not bool) {
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
