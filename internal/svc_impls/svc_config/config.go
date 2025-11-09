package svc_config

import (
	"context"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
)

var _ service.Config = &Service{}

// Service
// @autowire(service.Config,set=service)
type Service struct {
	Backends config.Backends
	Indexes  config.Indexes
}

func (s *Service) GetBackendList(ctx context.Context, req service.GetBackendListRequest) (resp service.GetBackendListResponse, err error) {
	for backendName, backendConfig := range s.Backends {
		resp.BackendList = append(resp.BackendList, service.Backend{
			Label: backendConfig.Label,
			Value: backendName,
		})
	}
	return
}

func (s *Service) GetStorageList(ctx context.Context, req service.GetStorageListRequest) (resp service.GetStorageListResponse, err error) {
	backendConfig, ok := s.Backends[req.BackendName]
	if !ok {
		return
	}
	for _, index := range backendConfig.Indexes {
		indexConfig := s.Indexes[index]
		resp.StorageList = append(resp.StorageList, service.Storage{
			Label: indexConfig.Label,
			Value: indexConfig.Name,
		})
	}
	return
}
