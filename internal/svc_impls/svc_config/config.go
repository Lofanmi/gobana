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
	BackendListConfig config.BackendList
}

func (s *Service) GetBackendList(ctx context.Context, req service.GetBackendListRequest) (resp service.GetBackendListResponse, err error) {
	for _, backend := range s.BackendListConfig {
		if !backend.Enabled {
			continue
		}
		resp.BackendList = append(resp.BackendList, service.Backend{
			Label: backend.Name,
			Value: backend.Name,
		})
	}
	return
}

func (s *Service) GetStorageList(ctx context.Context, req service.GetStorageListRequest) (resp service.GetStorageListResponse, err error) {
	var b *config.Backend
	for _, backend := range s.BackendListConfig {
		if !backend.Enabled {
			continue
		}
		if req.BackendName == backend.Name {
			b = &backend
			break
		}
	}
	if b == nil {
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
