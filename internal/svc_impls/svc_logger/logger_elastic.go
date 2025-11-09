package svc_logger

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/uritemplates"
)

func (s *Service) searchByElastic(ctx context.Context, backend config.BackendConfig, req service.SearchRequest) (resp service.SearchResponse, err error) {
	cli, err := s.ProviderFactory.GetProvider(backend.Name)
	if err != nil {
		return
	}

	m, rawQuery, err := s.elasticSearchResult(ctx, cli.ElasticClient, backend, req)
	if err != nil {
		return
	}
	resp.RawQuery = rawQuery
	resp.Count, resp.List, err = s.LogParser.ParseElastic(backend.Name, m)
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
		resp.Charts.Interval = req.ChartInterval
	}

	return
}

func (s *Service) elasticSearchResult(
	ctx context.Context,
	cli *elastic.Client,
	backend config.BackendConfig,
	req service.SearchRequest,
) (
	m map[string]*elastic.SearchResult,
	rawQuery map[string]any,
	err error,
) {
	queries, aggregations, err := s.QueryBuilder.BuildElastic(backend.Name, req)
	wg := new(sync.WaitGroup)
	mu := new(sync.RWMutex)
	m = map[string]*elastic.SearchResult{}
	rawQuery = map[string]any{}
	for _index, _query := range queries {
		_sortFields, ok := backend.SortFields[_index]
		if !ok {
			_sortFields = backend.SortFields[service.DefaultValue]
		}
		_aggregation := aggregations[_index]
		rawQuery[_index], _ = _query.Source()
		wg.Add(1)
		go func(index string, query elastic.Query, aggregation elastic.Aggregation, sortFields []config.SortField) {
			defer wg.Done()
			search := cli.Search()
			search.Index(s.Indexes[index].Meta.IndexMetaForElastic.Index).TrackTotalHits(req.TrackTotalHits).Query(query).Pretty(false).Version(true)
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

func (s *Service) searchDoElastic(ctx context.Context, backend config.BackendConfig, cli *elastic.Client, search *elastic.SearchService) (result *elastic.SearchResult, err error) {
	switch backend.Type {
	case service.ClientTypeElasticsearch:
		result, err = search.Do(ctx)
	case service.ClientTypeKibanaProxy:
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
