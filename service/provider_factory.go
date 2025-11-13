package service

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/olivere/elastic/v7"
)

type ProviderFactory interface {
	GetProvider(name string) (p ProviderClient, err error)
}

type ProviderClient struct {
	Type          ProviderType
	SlsClient     *sls.Client
	ElasticClient *elastic.Client
}
