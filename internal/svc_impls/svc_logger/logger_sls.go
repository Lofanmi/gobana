package svc_logger

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Lofanmi/gobana/service"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"golang.org/x/sync/errgroup"
)

var (
	reSelectCount = regexp.MustCompile(`(?i)select \*`)
)

func (s *Service) doSearchIndexSls(ctx context.Context, cli sls.ClientInterface, req service.SearchRequest, indexName, query string) (result service.SearchResultSls, err error) {
	result.RawQuery = query
	index := s.Indexes[indexName]
	project, store := index.Meta.IndexMetaForSls.Project, index.Meta.IndexMetaForSls.Store
	g, _ := errgroup.WithContext(ctx)
	from, to := req.TimeA/1000, req.TimeB/1000
	offset, limit := int64((req.PageNo-1)*req.PageSize), int64(req.PageSize)
	var searchStatement, analyticStatement string
	useAnalyticStatement := strings.Contains(query, "|")
	if useAnalyticStatement {
		pieces := strings.Split(query, "|")
		searchStatement, analyticStatement = strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1])
	} else {
		searchStatement = strings.TrimSpace(query)
	}

	slsRequest := &sls.GetLogRequest{From: from, To: to, Query: searchStatement, Lines: limit, Offset: offset, Reverse: true}
	if useAnalyticStatement {
		slsRequest.Query = fmt.Sprintf(`%s | %s ORDER BY __time__ DESC LIMIT %d, %d`, searchStatement, analyticStatement, offset, limit)
	}
	g.Go(func() (err error) {
		result.ResponseLog, result.ErrorByGetLogs = cli.GetLogsToCompletedV3(project, store, slsRequest)
		return
	})

	if !req.TrackTotalHits && !req.ChartVisible {
		_ = g.Wait()
		return
	}

	if useAnalyticStatement {
		if req.TrackTotalHits {
			g.Go(func() (err error) {
				slsRequestCount := &sls.GetLogRequest{From: from, To: to}
				slsRequestCount.Query = fmt.Sprintf(`%s | %s `, searchStatement, reSelectCount.ReplaceAllString(analyticStatement, "SELECT COUNT(*) as count"))
				result.ResponseCountByGetLogs, result.ErrorByGetLogs = cli.GetLogsToCompletedV3(project, store, slsRequestCount)
				return
			})
		}
		if req.ChartVisible {
			g.Go(func() (err error) {
				result.ResponseAggregation, result.ErrorResponseAggregation = cli.GetHistograms(project, store, "", from, to, slsRequest.Query)
				return
			})
		}
		_ = g.Wait()
		return
	}

	var (
		responseGetHistograms *sls.GetHistogramsResponse
		errorGetHistograms    error
	)
	g.Go(func() (err error) {
		responseGetHistograms, errorGetHistograms = cli.GetHistograms(project, store, "", from, to, slsRequest.Query)
		return
	})
	_ = g.Wait()

	if req.TrackTotalHits {
		result.ResponseCountByGetHistograms, result.ErrorByGetHistograms = responseGetHistograms, errorGetHistograms
	}
	if req.ChartVisible {
		result.ResponseAggregation, result.ErrorResponseAggregation = responseGetHistograms, errorGetHistograms
	}
	return
}
