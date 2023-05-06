package svc_logger

import (
	"context"
	"time"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/logic"
	"github.com/Lofanmi/gobana/service"
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
	_ = cli
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
