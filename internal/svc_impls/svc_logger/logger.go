package svc_logger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/constant"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/uritemplates"
)

var (
	_ service.Logger = &Service{}

	reSelectCount = regexp.MustCompile(`(?i)select \*`)
)

const (
	defaultPageNo        = 1
	defaultPageSize      = 20
	defaultMaxPageSize   = 200
	defaultMaxExportSize = 1000
)

// Service
// @autowire(service.Logger,set=service)
type Service struct {
	BackendListConfig config.BackendList
	BackendFactory    logic.BackendFactory
	QueryBuilder      logic.QueryBuilder
	LogParser         logic.LogParser
	AggregationParser logic.AggregationParser
}

func (s *Service) Search(ctx context.Context, req service.SearchRequest) (resp service.SearchResponse, err error) {
	if req.PageNo <= 0 {
		req.PageNo = defaultPageNo
	}
	if req.PageSize <= 0 {
		req.PageSize = defaultPageSize
	}
	if req.PageSize > defaultMaxPageSize {
		req.PageSize = defaultMaxPageSize
	}
	if req.TimeA == 0 || req.TimeB == 0 {
		t2 := time.Now()
		t1 := t2.Add(-time.Hour)
		req.TimeA, req.TimeB = t1.UnixMilli(), t2.UnixMilli()
	}
	if req.ChartVisible && req.ChartInterval <= 0 {
		req.ChartInterval = int32(defaultInterval(req.TimeA, req.TimeB))
	}
	defer func() {
		resp.PageNo = req.PageNo
		resp.PageSize = req.PageSize
		resp.TimeA = req.TimeA
		resp.TimeB = req.TimeB
	}()
	backend := s.BackendListConfig.Match(req.Backend)
	switch backend.Type {
	case constant.ClientTypeElasticsearch:
		fallthrough
	case constant.ClientTypeKibanaProxy:
		return s.searchByElastic(ctx, backend, req)
	case constant.ClientTypeSLS:
		return s.searchBySLS(ctx, backend, req)
	default:
		return
	}
}

func (s *Service) searchByElastic(ctx context.Context, backend config.Backend, req service.SearchRequest) (resp service.SearchResponse, err error) {
	cli, err := s.BackendFactory.GetBackendElastic(backend.Name)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(backend.Timeout))
	defer cancel()
	m, rawQuery, err := s.elasticSearchResult(ctx, cli, backend, req)
	if err != nil {
		return
	}
	resp.RawQuery = rawQuery
	resp.Count, resp.List, err = s.LogParser.ParseElastic(backend, m)
	if !req.TrackTotalHits {
		resp.Count = 10000
	}
	if req.ChartVisible {
		xAxis, yAxis, e := s.AggregationParser.ParseElastic(req.TimeA, req.TimeB, int64(req.ChartInterval), m)
		if e != nil {
			err = e
			return
		}
		resp.Charts.Legend = []string{"数量"}
		resp.Charts.XAxis = xAxis
		resp.Charts.Series.Name = "数量"
		resp.Charts.Series.Type = "bar"
		resp.Charts.Series.Symbol = "none"
		resp.Charts.Series.Smooth = true
		resp.Charts.Series.Data = yAxis
		resp.Charts.Interval = int(req.ChartInterval)
	}
	return
}

func (s *Service) elasticSearchResult(
	ctx context.Context,
	cli *elastic.Client,
	backend config.Backend,
	req service.SearchRequest,
) (
	m map[string]*elastic.SearchResult,
	rawQuery map[string]interface{},
	err error,
) {
	queries, aggregations, err := s.QueryBuilder.SearchQueryElastic(backend, req)
	wg := new(sync.WaitGroup)
	mu := new(sync.RWMutex)
	m = map[string]*elastic.SearchResult{}
	rawQuery = map[string]interface{}{}
	for _index, _query := range queries {
		_sortFields, ok := backend.SortFields[_index]
		if !ok {
			_sortFields = backend.SortFields[constant.DefaultValue]
		}
		_aggregation := aggregations[_index]
		rawQuery[_index], _ = _query.Source()
		wg.Add(1)
		go func(index string, query elastic.Query, aggregation elastic.Aggregation, sortFields []config.SortField) {
			defer wg.Done()
			search := cli.Search()
			search.Index(index).TrackTotalHits(req.TrackTotalHits).Query(query).Pretty(false).Version(true)
			search.From((req.PageNo - 1) * req.PageSize).Size(req.PageSize)
			if aggregation != nil {
				search.Aggregation("charts", aggregation)
			}
			for _, sortField := range sortFields {
				search.Sort(sortField.Field, sortField.Ascending)
			}
			result, e := s.searchDoElastic(ctx, backend, cli, search)
			if e == nil {
				mu.Lock()
				m[index] = result
				mu.Unlock()
			}
		}(_index, _query, _aggregation, _sortFields)
	}
	wg.Wait()
	return
}

func (s *Service) searchDoElastic(ctx context.Context, backend config.Backend, cli *elastic.Client, search *elastic.SearchService) (result *elastic.SearchResult, err error) {
	switch backend.Type {
	case constant.ClientTypeElasticsearch:
		result, err = search.Do(ctx)
	case constant.ClientTypeKibanaProxy:
		v := reflect.ValueOf(search).Elem()
		params := url.Values{}
		params.Set("method", http.MethodPost)
		field := v.FieldByName("index")
		var indexList []string
		for i := 0; i < field.Len(); i++ {
			indexList = append(indexList, field.Index(i).String())
		}
		if len(indexList) > 0 {
			path, _ := uritemplates.Expand("/{index}/_search", map[string]string{"index": strings.Join(indexList, ",")})
			params.Set("path", path)
		}
		field = v.FieldByName("searchSource")
		searchSource := (*elastic.SearchSource)(unsafe.Pointer(field.Pointer()))
		body, _ := searchSource.Source()
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")
		headers.Set("Cookie", backend.Auth.Cookie)
		headers.Set("kbn-version", backend.Auth.KbnVersion)
		var res *elastic.Response
		if res, err = cli.PerformRequest(ctx, elastic.PerformRequestOptions{
			Method:  http.MethodPost,
			Path:    "/api/console/proxy",
			Params:  params,
			Body:    body,
			Headers: headers,
		}); err != nil {
			return
		}
		if err = json.Unmarshal(res.Body, &result); err != nil {
			result.Header = res.Header
			return
		}
		result.Header = res.Header
	}
	return
}

func (s *Service) searchBySLS(ctx context.Context, backend config.Backend, req service.SearchRequest) (resp service.SearchResponse, err error) {
	cli, err := s.BackendFactory.GetBackendSLS(backend.Name)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(backend.Timeout))
	defer cancel()
	m, rawQuery, err := s.slsSearchResult(ctx, cli, backend, req)
	_ = m
	if err != nil {
		return
	}
	resp.RawQuery = rawQuery
	resp.Count, resp.List, err = s.LogParser.ParseSLS(backend, m)
	if req.ChartVisible {
		xAxis, yAxis, e := s.AggregationParser.ParseSLS(req.TimeA, req.TimeB, int64(req.ChartInterval), m)
		if e != nil {
			err = e
			return
		}
		resp.Charts.Legend = []string{"数量"}
		resp.Charts.XAxis = xAxis
		resp.Charts.Series.Name = "数量"
		resp.Charts.Series.Type = "bar"
		resp.Charts.Series.Symbol = "none"
		resp.Charts.Series.Smooth = true
		resp.Charts.Series.Data = yAxis
		resp.Charts.Interval = int(req.ChartInterval)
	}
	return
}

func (s *Service) slsSearchResult(
	ctx context.Context,
	cli sls.ClientInterface,
	backend config.Backend,
	req service.SearchRequest,
) (
	m map[string]logic.SLSSearchResult,
	rawQuery map[string]interface{},
	err error,
) {
	queries, err := s.QueryBuilder.SearchQuerySLS(backend, req)
	wg := new(sync.WaitGroup)
	mu := new(sync.RWMutex)
	m = map[string]logic.SLSSearchResult{}
	rawQuery = map[string]interface{}{}
	for _index, _query := range queries {
		rawQuery[_index] = _query
		wg.Add(1)
		go func(index, query string) {
			defer wg.Done()
			result, e := s.searchDoSLS(ctx, backend, cli, req, index, query)
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

func (s *Service) searchDoSLS(
	ctx context.Context,
	backend config.Backend,
	cli sls.ClientInterface,
	req service.SearchRequest,
	index, query string,
) (result logic.SLSSearchResult, err error) {
	_ = ctx
	pieces := strings.Split(index, "|")
	if len(pieces) != 2 {
		return
	}
	project, store := pieces[0], pieces[1]
	wg := new(sync.WaitGroup)
	mu := new(sync.RWMutex)
	from, to := req.TimeA/1000, req.TimeB/1000
	offset, limit := int64((req.PageNo-1)*req.PageSize), int64(req.PageSize)

	var (
		phrase, where string
	)
	if strings.Contains(query, "|") {
		pieces = strings.Split(query, "|")
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

	switch backend.Type {
	case constant.ClientTypeSLS:
		if req.TrackTotalHits {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if where != "" {
					slsRequestCount := &sls.GetLogRequest{From: from, To: to}
					slsRequestCount.Query = fmt.Sprintf(`%s | %s `, phrase, reSelectCount.ReplaceAllString(where, "SELECT COUNT(*) as count"))
					resp, e := cli.GetLogsV3(project, store, slsRequestCount)
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
			resp, e := cli.GetLogsV3(project, store, slsRequest)
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
	}
	wg.Wait()
	return
}

func defaultInterval(timeA, timeB int64) (interval int) {
	interval = (int)((timeB - timeA) / 1000 / constant.MaxChartPoints)
	if interval <= 1 {
		interval = 1
	} else if interval <= 5 {
		interval = 5
	} else if interval <= 10 {
		interval = 10
	} else if interval <= 30 {
		interval = 30
	} else if interval <= 60 {
		interval = 60
	} else if interval <= 300 {
		interval = 300
	} else if interval <= 900 {
		interval = 900
	} else if interval <= 1800 {
		interval = 1800
	} else if interval <= 3600 {
		interval = 3600
	} else if interval <= 3600*3 {
		interval = 3600 * 3
	} else if interval <= 3600*9 {
		interval = 3600 * 9
	} else if interval <= 3600*12 {
		interval = 3600 * 12
	} else {
		interval = 3600 * 24
	}
	return
}
