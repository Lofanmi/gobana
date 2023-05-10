package svc_logger

import (
	"context"
	"sync"
	"time"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
	"github.com/olivere/elastic/v7"
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
	case logic.ClientTypeElasticsearch:
		fallthrough
	case logic.ClientTypeKibanaProxy:
		return s.searchByElastic(ctx, backend, req)
	case logic.ClientTypeSLS:
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
	for i, q := range queries {
		rawQuery[i], _ = q.Source()
		wg.Add(1)
		go func(index string, query elastic.Query, sortFields []config.SortField) {
			defer wg.Done()
			search := cli.Search()
			search.Index(index).TrackTotalHits(trackTotalHits).Query(query).Pretty(false).Version(true)
			search.From((req.PageNo - 1) * req.PageSize).Size(req.PageSize)
			for _, sortField := range sortFields {
				search.Sort(sortField.Field, sortField.Ascending)
			}
			result, e := search.Do(ctx)
			if e == nil {
				mu.Lock()
				m[index] = result
				mu.Unlock()
			}
		}(i, q, backend.SortFields[i])
	}
	wg.Wait()
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
