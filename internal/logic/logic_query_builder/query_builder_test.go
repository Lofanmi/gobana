package logic_query_builder

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/constant"
	"github.com/Lofanmi/gobana/service"
)

func TestQueryBuilder_SearchQueryElastic(t *testing.T) {
	const indexName = "index-json-log"
	t2 := time.Now()
	t1 := t2.Add(-time.Hour)
	type args struct {
		backend config.Backend
		req     service.SearchRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				backend: config.Backend{
					Type:          constant.ClientTypeElasticsearch,
					MultiSearch:   map[string]config.MultiSearch{"程序日志": {IndexList: []string{indexName}}},
					DefaultFields: map[string][]string{indexName: {"host", "tag"}},
					BuildInQueries: map[string]config.BuildInQuery{
						indexName: {
							Must: []config.BuildInQueryEntry{
								{
									Name:     "必看业务线",
									Field:    "app_name",
									Values:   []interface{}{"app_1", "app_2"},
									Operator: "and",
									Always:   true,
								},
							},
							MustNot: []config.BuildInQueryEntry{
								{
									Name:     "不看机器",
									Field:    "host",
									Values:   []interface{}{"host_1", "host_2"},
									Operator: "and",
									Always:   true,
								},
							},
							Or: []config.BuildInQueryEntry{
								{
									Name:     "都可以",
									Field:    "tag",
									Values:   []interface{}{"tag_1", "tag_2"},
									Operator: "and",
									Always:   true,
								},
							},
						},
					},
				},
				req: service.SearchRequest{
					TimeA:   t1.UnixMilli(),
					TimeB:   t2.UnixMilli(),
					Backend: "",
					Storage: "程序日志",
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
			s := &QueryBuilder{}
			gotQueries, _, err := s.SearchQueryElastic(tt.args.backend, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchQueryElastic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for index, query := range gotQueries {
				m, _ := query.Source()
				data, _ := json.MarshalIndent(&m, "", "    ")
				t.Logf("%s:\n%s", index, string(data))
			}
		})
	}
}
