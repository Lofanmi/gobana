package svc_config

import (
	"context"
	"sort"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
)

var _ service.Config = &Service{}

// Service
// @autowire(service.Config,set=service)
type Service struct {
	BackendListConfig config.Backends
}

func (s *Service) GetBackendList(ctx context.Context, req service.GetBackendListRequest) (resp service.GetBackendListResponse, err error) {
	for backendName, backendConfig := range s.BackendListConfig {
		resp.BackendList = append(resp.BackendList, service.Backend{
			Label: backendConfig.Label,
			Value: backendName,
		})
	}
	return
}

func (s *Service) GetStorageList(ctx context.Context, req service.GetStorageListRequest) (resp service.GetStorageListResponse, err error) {
	backendConfig, ok := s.BackendListConfig[req.BackendName]
	if !ok {
		return
	}

	var multiSearchList config.MultiSearchSlice
	for _, multiSearch := range b.MultiSearch {
		multiSearchList = append(multiSearchList, multiSearch)
	}
	sort.Sort(multiSearchList)
	for _, multiSearch := range multiSearchList {
		resp.StorageList = append(resp.StorageList, service.Storage{
			Label: multiSearch.Name,
			Value: multiSearch.Name,
		})
	}
	return
}
