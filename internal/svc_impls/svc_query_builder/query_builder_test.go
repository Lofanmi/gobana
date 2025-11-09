package svc_query_builder

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
)

const (
	backendName = "access_log"
	indexName   = "index-json-log"
)

var (
	backend = config.BackendConfig{
		Name:          backendName,
		Indexes:       []string{indexName},
		DefaultFields: map[string][]string{indexName: {"host", "tag"}},
		BuildInQuery: map[string]config.BuildInQuery{
			indexName: {
				Must: []config.BuildInQueryEntry{
					{
						Name:     "必看业务线",
						Field:    "app_name",
						Values:   []any{"app_1", "app_2"},
						Operator: "and",
						Always:   true,
					},
				},
				MustNot: []config.BuildInQueryEntry{
					{
						Name:     "不看机器",
						Field:    "host",
						Values:   []any{"host_1", "host_2"},
						Operator: "and",
						Always:   true,
					},
				},
				Or: []config.BuildInQueryEntry{
					{
						Name:     "都可以",
						Field:    "tag",
						Values:   []any{"tag_1", "tag_2"},
						Operator: "and",
						Always:   true,
					},
				},
			},
		},
	}
	backends = config.Backends{backendName: backend}
)

func TestQueryBuilder_BuildElastic_QueryTypeByHuman(t *testing.T) {
	t2 := time.Now()
	t1 := t2.Add(-time.Hour)
	type args struct {
		backendName string
		req         service.SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				backendName: backendName,
				req: service.SearchRequest{
					TimeA:   t1.UnixMilli(),
					TimeB:   t2.UnixMilli(),
					QueryBy: service.QueryTypeByHuman,
					Query: service.QueryByHuman{
						Must:    []string{"1.2.3.4", "5.6.7.8"},
						MustNot: []string{"127.0.0.1", "192.168.1.1"},
						Or:      []string{"8.8.8.8", "9.9.9.9"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &QueryBuilder{Backends: backends}
			gotQueries, gotAggregations, err := s.BuildElastic(tt.args.backendName, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildElastic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for index, query := range gotQueries {
				m, _ := query.Source()
				data, _ := json.MarshalIndent(&m, "", "    ")
				t.Logf("%s:\n%s", index, string(data))
				if v := gotAggregations[index]; v != nil {
					m2, _ := v.Source()
					data2, _ := json.MarshalIndent(&m2, "", "    ")
					t.Logf("%s:\n%s", index, string(data2))
				}
			}
		})
	}
}

func TestQueryBuilder_BuildElastic_QueryTypeByLucene(t *testing.T) {
	t2 := time.Now()
	t1 := t2.Add(-time.Hour)
	type args struct {
		backendName string
		req         service.SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				backendName: backendName,
				req: service.SearchRequest{
					TimeA:   t1.UnixMilli(),
					TimeB:   t2.UnixMilli(),
					Backend: "",
					Storage: "程序日志",
					QueryBy: service.QueryTypeByLucene,
					Query: service.QueryByLucene{
						Lucene: "status:200 AND message:client",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &QueryBuilder{}
			gotQueries, gotAggregations, err := s.BuildElastic(tt.args.backendName, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchQueryElastic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for index, query := range gotQueries {
				m, _ := query.Source()
				data, _ := json.MarshalIndent(&m, "", "    ")
				t.Logf("%s:\n%s", index, string(data))
				if v := gotAggregations[index]; v != nil {
					m2, _ := v.Source()
					data2, _ := json.MarshalIndent(&m2, "", "    ")
					t.Logf("%s:\n%s", index, string(data2))
				}
			}
		})
	}
}

func TestSlsMainQuery(t *testing.T) {
	var mainQuery *SlsQuery

	mainQuery = new(SlsQuery)
	t.Log(mainQuery)

	fields := []string{"application", "client_ip", "content"}
	conditions := []string{"%gaia%", "%timer%", "sdk", "login"}

	mainQuery = new(SlsQuery)
	searchConditions, fuzzyConditions := quoteConditions(conditions)
	mainQuery.PrepareSearchConditions(searchConditions, operatorAnd, false)
	mainQuery.PrepareFuzzyConditions(fields, fuzzyConditions, operatorAnd, false)
	t.Log(mainQuery)
	mainQuery.PrepareSearchConditions([]string{"'abc'"}, operatorAnd, false)
	mainQuery.PrepareFuzzyConditions([]string{"host"}, []string{"'%test_host%'"}, operatorAnd, false)
	t.Log(mainQuery)

	mainQuery = new(SlsQuery)
	searchConditions, fuzzyConditions = quoteConditions(conditions)
	mainQuery.PrepareSearchConditions(searchConditions, operatorOr, false)
	mainQuery.PrepareFuzzyConditions(fields, fuzzyConditions, operatorOr, false)
	t.Log(mainQuery)

	mainQuery = new(SlsQuery)
	searchConditions, fuzzyConditions = quoteConditions(conditions)
	mainQuery.PrepareSearchConditions(searchConditions, operatorNot, true)
	mainQuery.PrepareFuzzyConditions(fields, fuzzyConditions, operatorNot, true)
	t.Log(mainQuery)
}
