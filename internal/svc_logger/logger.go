package svc_logger

import (
	"context"

	"github.com/Lofanmi/gobana/service"
)

var _ service.Logger = &Service{}

// Service
// @autowire(service.Logger,set=service)
type Service struct {
}

// NewService 日志服务
func NewService() *Service {
	return &Service{}
}

func (s Service) Search(ctx context.Context, req service.SearchRequest) (resp service.SearchResponse, err error) {
	return
}
