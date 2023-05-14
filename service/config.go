package service

import (
	"context"
)

type Config interface {
	GetBackendList(ctx context.Context, req GetBackendListRequest) (resp GetBackendListResponse, err error)
	GetStorageList(ctx context.Context, req GetStorageListRequest) (resp GetStorageListResponse, err error)
}

type GetBackendListRequest struct{}

type GetBackendListResponse struct {
	BackendList []Backend `json:"backend_list"`
}

type GetStorageListRequest struct {
	BackendName string `json:"backend_name" form:"backend_name" binding:"required"`
}

type GetStorageListResponse struct {
	StorageList []Storage `json:"storage_list"`
}

type Backend struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Storage struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
