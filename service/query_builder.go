package service

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/olivere/elastic/v7"
)

type QueryBuilder interface {
	BuildSls(indexName string, req SearchRequest) (queryRes string, err error)
	BuildElastic(indexName string, req SearchRequest) (queryRes elastic.Query, aggregationRes elastic.Aggregation, err error)
}

type SearchResult struct {
	ProviderType        ProviderType        `json:"provider_type"`
	SearchResultSls     SearchResultSls     `json:"sls"`
	SearchResultElastic SearchResultElastic `json:"elastic"`
}

type SearchResultSls struct {
	RawQuery                     string                     `json:"raw_query"`
	ResponseCountByGetHistograms *sls.GetHistogramsResponse `json:"response_count_by_get_histograms"`
	ErrorByGetHistograms         error                      `json:"error_by_get_histograms"`
	ResponseCountByGetLogs       *sls.GetLogsV3Response     `json:"response_count_by_get_logs"`
	ErrorByGetLogs               error                      `json:"error_by_get_logs"`
	ResponseLog                  *sls.GetLogsV3Response     `json:"response_log"`
	ErrorResponseLog             error                      `json:"error_response_log"`
	ResponseAggregation          *sls.GetHistogramsResponse `json:"response_aggregation"`
	ErrorResponseAggregation     error                      `json:"error_response_aggregation"`
}

type SearchResultElastic struct {
	RawQuery      any                   `json:"raw_query"`
	Response      *elastic.SearchResult `json:"response"`
	ErrorResponse error                 `json:"error_response"`
}
