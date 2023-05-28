package logic

import (
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/olivere/elastic/v7"
)

type SLSSearchResult struct {
	ResponseCount       *sls.GetLogsResponse       `json:"response_count"`
	ResponseLog         *sls.GetLogsResponse       `json:"response_log"`
	ResponseAggregation *sls.GetHistogramsResponse `json:"response_aggregation"`
}

type LogParser interface {
	ParseElastic(backend config.Backend, m map[string]*elastic.SearchResult) (total int, logs service.LogItems, err error)
	ParseSLS(backend config.Backend, m map[string]SLSSearchResult) (total int, logs service.LogItems, err error)
}
