package service

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/olivere/elastic/v7"
)

type SlsSearchResult struct {
	ResponseCountByGetHistograms *sls.GetHistogramsResponse `json:"response_count_by_get_histograms"`
	ResponseCountByGetLogs       *sls.GetLogsV3Response     `json:"response_count_by_get_logs"`
	ResponseLog                  *sls.GetLogsV3Response     `json:"response_log"`
	ResponseAggregation          *sls.GetHistogramsResponse `json:"response_aggregation"`
	ErrorByGetHistograms         error                      `json:"error_by_get_histograms"`
	ErrorByGetLogs               error                      `json:"error_by_get_logs"`
	ErrorResponseLog             error                      `json:"error_response_log"`
	ErrorResponseAggregation     error                      `json:"error_response_aggregation"`
}

type LogParser interface {
	ParseElastic(backendName string, m map[string]*elastic.SearchResult) (total int, logs LogItems, err error)
	ParseSLS(backendName string, m map[string]SlsSearchResult) (total int, logs LogItems, err error)
}
