package svc_logger

import (
	"context"
	"regexp"
	"time"

	"github.com/Lofanmi/cmap"
	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/gotil"
	"github.com/Lofanmi/gobana/service"
	"golang.org/x/sync/errgroup"
)

var (
	_ service.Logger = &Service{}

	reSelectCount = regexp.MustCompile(`(?i)select \*`)
)

const (
	defaultPageNo      = 1
	defaultPageSize    = 20
	defaultMaxPageSize = 200
)

// Service
// @autowire(service.Logger,set=service)
type Service struct {
	Indexes         config.Indexes
	Backends        config.Backends
	Providers       config.Providers
	ProviderFactory service.ProviderFactory
	QueryBuilder    service.QueryBuilder
	GoJa            service.GoJa
}

func (s *Service) Search(ctx context.Context, req service.SearchRequest) (resp service.SearchResponse, err error) {
	req.PageNo = gotil.IfThen(req.PageNo <= 0, defaultPageNo)
	req.PageSize = gotil.IfThen(req.PageSize <= 0, defaultPageSize)
	req.PageSize = gotil.IfThen(req.PageNo > defaultMaxPageSize, defaultPageNo)
	if req.TimeA == 0 || req.TimeB == 0 {
		t2 := time.Now()
		t1 := t2.Add(-time.Hour)
		req.TimeA, req.TimeB = t1.UnixMilli(), t2.UnixMilli()
	}
	if req.ChartVisible && req.ChartInterval <= 0 {
		req.ChartInterval = gotil.DefaultInterval(req.TimeA, req.TimeB, service.MaxChartPoints)
	}
	defer func() {
		resp.PageNo, resp.PageSize = req.PageNo, req.PageSize
		resp.TimeA, resp.TimeB = req.TimeA, req.TimeB
	}()

	return
}

func (s *Service) doSearchIndexes(ctx context.Context, req service.SearchRequest, resp *service.SearchResponse) (err error) {
	backend := s.Backends[req.Backend]
	m := cmap.NewStringHashMap[service.SearchResult](cmap.WithShardCount(2))
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < len(backend.Indexes); i++ {
		if err = s.doSearchIndex(ctx, backend.Indexes[i], req, m, g); err != nil {
			return
		}
	}
	if err = g.Wait(); err != nil {
		return
	}

	for _, indexName := range backend.Indexes {
		searchResult, _ := m.Get(indexName)
		switch searchResult.ProviderType {
		case service.ProviderTypeSls:
		case service.ProviderTypeElasticsearch:
		}
	}

	// resp.RawQuery = result.RawQuery
	// resp.Count, resp.List, err = s.LogParser.ParseElastic("", nil)

	if req.TrackTotalHits {
		resp.Count = s.ParseTotal(m)
	} else {
		resp.Count = 10000
	}

	if !req.ChartVisible {
		return
	}
	xAxis, yAxis, err := s.ParseAggregation(req.TimeA, req.TimeB, int64(req.ChartInterval), m)
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

func (s *Service) doSearchIndex(ctx context.Context, indexName string, req service.SearchRequest, res *cmap.Map[string, service.SearchResult], g *errgroup.Group) (err error) {
	index, ok := s.Indexes[indexName]
	if !ok {
		return
	}
	p, err := s.ProviderFactory.GetProvider(index.ProviderName)
	if err != nil {
		return
	}
	g.Go(func() error {
		switch p.Type {
		case service.ProviderTypeSls:
			query, e := s.QueryBuilder.BuildSls(indexName, req)
			if e != nil {
				return e
			}
			result, e := s.doSearchIndexSls(ctx, p.SlsClient, req, indexName, query)
			if e != nil {
				return e
			}
			if result.ErrorByGetHistograms != nil {
				return result.ErrorByGetHistograms
			}
			if result.ErrorByGetLogs != nil {
				return result.ErrorByGetLogs
			}
			if result.ErrorResponseLog != nil {
				return result.ErrorResponseLog
			}
			if result.ErrorResponseAggregation != nil {
				return result.ErrorResponseAggregation
			}
			res.Put(indexName, service.SearchResult{ProviderType: p.Type, SearchResultSls: result})
		case service.ProviderTypeElasticsearch:
			query, aggregation, e := s.QueryBuilder.BuildElastic(indexName, req)
			if e != nil {
				return e
			}
			result, e := s.doSearchIndexElastic(ctx, p.ElasticClient, req, indexName, query, aggregation)
			if e != nil {
				return e
			}
			if result.ErrorResponse != nil {
				return result.ErrorResponse
			}
			res.Put(indexName, service.SearchResult{ProviderType: p.Type, SearchResultElastic: result})
		}
		return nil
	})
	return
}
