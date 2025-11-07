package config

import (
	"github.com/spf13/cast"
)

// ---------------------------------------------------------------------------------------------------------------------

type BuildInQuery struct {
	Must    BuildInQueryEntrySlice `json:"must"`
	MustNot BuildInQueryEntrySlice `json:"must_not"`
	Or      BuildInQueryEntrySlice `json:"or"`
}

type BuildInQueryEntrySlice []BuildInQueryEntry

type BuildInQueryEntry struct {
	Name     string `json:"name"`
	Field    string `json:"field"`
	Values   []any  `json:"values"`
	Operator string `json:"operator"`
	Always   bool   `json:"always"`
}

func (s BuildInQueryEntrySlice) FieldConditions(fn func(field string, conditions []string) bool) {
	for _, item := range s {
		if !item.Always {
			continue
		}
		if !fn(item.Field, cast.ToStringSlice(item.Values)) {
			break
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------
