package logic_backend_factory

import (
	"errors"
	"net/http"
	"sync"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/internal/constant"
	"github.com/Lofanmi/gobana/internal/logic"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/olivere/elastic/v7"
)

var (
	_ logic.BackendFactory = &BackendFactory{}
)

// BackendFactory
// @autowire(set=logics)
type BackendFactory struct {
	BackendListConfig config.BackendList
	m                 map[string]interface{}
	mu                *sync.Mutex
	httpClient        *http.Client
}

func NewBackendFactory(backendListConfig config.BackendList) logic.BackendFactory {
	s := &BackendFactory{
		BackendListConfig: backendListConfig,
		m:                 map[string]interface{}{},
		mu:                new(sync.Mutex),
		httpClient:        new(http.Client),
	}
	return s
}

func (s *BackendFactory) GetBackendElastic(name string) (cli *elastic.Client, err error) {
	res, err := s.GetBackend(name)
	if err != nil {
		return
	}
	if v, ok := res.(*elastic.Client); !ok {
		err = errors.New("GetBackend failed")
	} else {
		cli = v
	}
	return
}

func (s *BackendFactory) GetBackendSLS(name string) (cli sls.ClientInterface, err error) {
	res, err := s.GetBackend(name)
	if err != nil {
		return
	}
	if v, ok := res.(sls.ClientInterface); !ok {
		err = errors.New("GetBackendSLS failed")
	} else {
		cli = v
	}
	return
}

func (s *BackendFactory) GetBackend(name string) (res interface{}, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, backend := range s.BackendListConfig {
		if backend.Name != name {
			continue
		}
		return s.getClient(backend)
	}
	return
}

func (s *BackendFactory) getClient(config config.Backend) (res interface{}, err error) {
	typ := config.Type
	if _, ok := s.m[typ]; !ok {
		switch typ {
		case constant.ClientTypeElasticsearch:
			fallthrough
		case constant.ClientTypeKibanaProxy:
			var cli *elastic.Client
			cli, err = elastic.NewClient(
				elastic.SetURL(config.Addr),
				elastic.SetHttpClient(s.httpClient),
				elastic.SetHealthcheck(false),
				elastic.SetSniff(false),
			)
			if err != nil {
				return
			}
			s.m[typ] = cli
		case constant.ClientTypeSLS:
			cli := &sls.Client{
				Endpoint:        config.Addr,
				AccessKeyID:     config.Auth.AccessKeyID,
				AccessKeySecret: config.Auth.AccessKeySecret,
				SecurityToken:   "",
				HTTPClient:      s.httpClient,
			}
			s.m[typ] = cli
		}
	}
	res = s.m[typ]
	return
}
