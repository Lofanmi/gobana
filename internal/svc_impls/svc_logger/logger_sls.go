package svc_logger

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	sls "github.com/aliyun/aliyun-log-go-sdk"
)

func (s *Service) searchBySLS(ctx context.Context, backend config.BackendConfig, req service.SearchRequest) (resp service.SearchResponse, err error) {
	cli, err := s.ProviderFactory.GetProvider(backend.Name)
	if err != nil {
		return
	}
	m, rawQuery, err := s.slsSearchResult(ctx, cli.SlsClient, backend, req)
	if err != nil {
		return
	}
	resp.RawQuery = rawQuery
	resp.Count, resp.List, err = s.LogParser.ParseSLS(backend.Name, m)
	if !req.ChartVisible {
		return
	}
	xAxis, yAxis, err := s.AggregationParser.ParseSLS(req.TimeA, req.TimeB, int64(req.ChartInterval), m)
	if err != nil {
		return
	}
	resp.Charts.Legend = []string{"数量"}
	resp.Charts.XAxis = xAxis
	resp.Charts.Series.Name = "数量"
	resp.Charts.Series.Type = "bar"
	resp.Charts.Series.Symbol = "none"
	resp.Charts.Series.Smooth = true
	resp.Charts.Series.Data = yAxis
	resp.Charts.Interval = req.ChartInterval
	return
}

func (s *Service) slsSearchResult(ctx context.Context, cli sls.ClientInterface, backend config.BackendConfig, req service.SearchRequest) (
	m map[string]service.SlsSearchResult,
	rawQuery map[string]any,
	err error,
) {
	queries, err := s.QueryBuilder.BuildSls(backend.Name, req)
	wg := new(sync.WaitGroup)
	mu := new(sync.RWMutex)
	m = map[string]service.SlsSearchResult{}
	rawQuery = map[string]any{}
	for _index, _query := range queries {
		rawQuery[_index] = _query
		wg.Add(1)
		go func(index, query string) {
			defer wg.Done()
			result, e := s.searchDoSLS(ctx, cli, req, project, store, query)
			if e == nil {
				mu.Lock()
				m[index] = result
				mu.Unlock()
			}
		}(_index, _query)
	}
	wg.Wait()
	return
}

func (s *Service) provider(_ context.Context, indexName string, req service.SearchRequest, query string, fn func()) (err error) {
	index, ok := s.Indexes[indexName]
	if !ok {
		return
	}
	p, err := s.ProviderFactory.GetProvider(index.ProviderName)
	if err != nil {
		return
	}
	switch p.Type {
	case service.ProviderTypeSls:
		project, store := index.Meta.IndexMetaForSls.Project, index.Meta.IndexMetaForSls.Store

	case service.ProviderTypeElasticsearch:
	}
	return
}

func (s *Service) searchDoSLS(_ context.Context, cli sls.ClientInterface, req service.SearchRequest, project, store, query string) (result service.SlsSearchResult, err error) {
	wg := new(sync.WaitGroup)
	mu := new(sync.RWMutex)
	from, to := req.TimeA/1000, req.TimeB/1000
	offset, limit := int64((req.PageNo-1)*req.PageSize), int64(req.PageSize)
	phrase, where := "", ""
	if strings.Contains(query, "|") {
		pieces := strings.Split(query, "|")
		phrase, where = strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1])
	} else {
		phrase = strings.TrimSpace(query)
	}

	slsRequest := &sls.GetLogRequest{From: from, To: to, Query: query, Lines: limit, Offset: offset, Reverse: true}
	if where != "" {
		slsRequest.Query = fmt.Sprintf(`%s | %s ORDER BY __time__ DESC LIMIT %d, %d`, phrase, where, offset, limit)
	} else {
		slsRequest.Query = phrase
	}

	if req.TrackTotalHits {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if where != "" {
				slsRequestCount := &sls.GetLogRequest{From: from, To: to}
				slsRequestCount.Query = fmt.Sprintf(`%s | %s `, phrase, reSelectCount.ReplaceAllString(where, "SELECT COUNT(*) as count"))
				resp, e := cli.GetLogsToCompletedV3(project, store, slsRequestCount)
				if e != nil {
					mu.Lock()
					result.ErrorByGetLogs = e
					mu.Unlock()
					return
				}
				mu.Lock()
				result.ResponseCountByGetLogs = resp
				mu.Unlock()
				return
			}
			resp, e := cli.GetHistograms(project, store, "", from, to, slsRequest.Query)
			if e != nil {
				mu.Lock()
				result.ErrorByGetHistograms = e
				mu.Unlock()
				return
			}
			mu.Lock()
			result.ResponseCountByGetHistograms = resp
			mu.Unlock()
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, e := cli.GetLogsToCompletedV3(project, store, slsRequest)
		if e != nil {
			mu.Lock()
			result.ErrorByGetLogs = e
			mu.Unlock()
			return
		}
		mu.Lock()
		result.ResponseLog = resp
		mu.Unlock()
	}()
	if req.ChartVisible {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, e := cli.GetHistograms(project, store, "", from, to, slsRequest.Query)
			if e != nil {
				mu.Lock()
				result.ErrorResponseAggregation = e
				mu.Unlock()
				return
			}
			mu.Lock()
			result.ResponseAggregation = resp
			mu.Unlock()
		}()
	}
	wg.Wait()
	return
}
