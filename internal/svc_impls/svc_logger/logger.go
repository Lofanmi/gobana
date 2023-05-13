package svc_logger

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/constant"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/uritemplates"
)

var _ service.Logger = &Service{}

const (
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
}

func (s *Service) Search(ctx context.Context, req service.SearchRequest) (resp service.SearchResponse, err error) {
	if req.PageSize <= 0 || req.PageSize > defaultMaxPageSize {
		req.PageSize = defaultMaxPageSize
	}
	if req.TimeA == 0 || req.TimeB == 0 {
		t2 := time.Now()
		t1 := t2.Add(-time.Hour)
		req.TimeA, req.TimeB = t1.UnixMilli(), t2.UnixMilli()
	}
	resp.PageNo = req.PageNo
	resp.PageSize = req.PageSize
	resp.TimeA = req.TimeA
	resp.TimeB = req.TimeB

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
	m, rawQuery, trackTotalHits, err := s.elasticSearchResult(ctx, cli, backend, req)
	if err != nil {
		return
	}
	resp.RawQuery = rawQuery
	resp.Count, resp.List, err = s.LogParser.ParseElastic(backend, m)
	if !trackTotalHits {
		resp.Count = 10000
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
	trackTotalHits bool,
	err error,
) {
	queries, trackTotalHits, err := s.QueryBuilder.SearchQueryElastic(backend, req)
	wg := new(sync.WaitGroup)
	mu := new(sync.RWMutex)
	m = map[string]*elastic.SearchResult{}
	rawQuery = map[string]interface{}{}
	for _index, _query := range queries {
		_sortFields, ok := backend.SortFields[_index]
		if !ok {
			_sortFields = backend.SortFields[constant.DefaultValue]
		}
		rawQuery[_index], _ = _query.Source()
		wg.Add(1)
		go func(index string, query elastic.Query, sortFields []config.SortField) {
			defer wg.Done()
			search := cli.Search()
			search.Index(index).TrackTotalHits(trackTotalHits).Query(query).Pretty(false).Version(true)
			search.From((req.PageNo - 1) * req.PageSize).Size(req.PageSize)
			for _, sortField := range sortFields {
				search.Sort(sortField.Field, sortField.Ascending)
			}
			result, e := s.searchDo(ctx, backend, cli, search)
			if e == nil {
				mu.Lock()
				m[index] = result
				mu.Unlock()
			}
		}(_index, _query, _sortFields)
	}
	wg.Wait()
	return
}

func (s *Service) searchDo(ctx context.Context, backend config.Backend, cli *elastic.Client, search *elastic.SearchService) (result *elastic.SearchResult, err error) {
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
	_ = cli
	return
}
