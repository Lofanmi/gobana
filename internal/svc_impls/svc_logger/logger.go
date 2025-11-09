package svc_logger

import (
	"context"
	"regexp"
	"time"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
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
	Indexes           config.Indexes
	Backends          config.Backends
	ProviderFactory   service.ProviderFactory
	QueryBuilder      service.QueryBuilder
	LogParser         service.LogParser
	AggregationParser service.AggregationParser
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
		req.ChartInterval = defaultInterval(req.TimeA, req.TimeB)
	}
	defer func() {
		resp.PageNo = req.PageNo
		resp.PageSize = req.PageSize
		resp.TimeA = req.TimeA
		resp.TimeB = req.TimeB
	}()
	backend := s.Backends[req.Backend]
	switch backend.Type {
	case service.ClientTypeElasticsearch:
		fallthrough
	case service.ClientTypeKibanaProxy:
		return s.searchByElastic(ctx, backend, req)
	case service.ClientTypeSLS:
		return s.searchBySLS(ctx, backend, req)
	default:
		return
	}
}

func defaultInterval(timeA, timeB int64) (interval int) {
	interval = (int)((timeB - timeA) / 1000 / service.MaxChartPoints)
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
