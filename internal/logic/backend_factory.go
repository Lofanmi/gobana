package logic

import (
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/olivere/elastic/v7"
)

type SLSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SecurityToken   string `json:"security_token"`
}

type BackendFactory interface {
	GetBackendElastic(name string) (cli *elastic.Client, err error)
	GetBackendSLS(name string) (cli sls.ClientInterface, err error)
	GetBackend(name string) (res interface{}, err error)
}
