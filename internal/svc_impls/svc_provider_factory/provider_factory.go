package svc_provider_factory

import (
	"errors"
	"net/http"
	"time"

	"github.com/Lofanmi/gobana/internal/config"
	"github.com/Lofanmi/gobana/service"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/olivere/elastic/v7"
)

var (
	_ service.ProviderFactory = &ProviderFactory{}
)

// ProviderFactory
// @autowire(service.ProviderFactory,set=service)
type ProviderFactory struct {
	Providers config.Providers
}

func (s *ProviderFactory) GetProvider(name string) (p service.ProviderClient, err error) {
	cfg := s.Providers.Match(name)
	if cfg.Name == "" {
		err = errors.New("no such provider")
		return
	}
	p.Type = cfg.Type

	httpClient := new(http.Client)
	httpClient.Timeout = time.Duration(cfg.TimeoutMillisecond) * time.Millisecond

	switch cfg.Type {
	case service.ProviderTypeSls:
		p.SlsClient, err = s.getProviderSLS(&cfg, httpClient)
	case service.ProviderTypeElasticsearch:
		p.ElasticClient, err = s.getProviderElastic(&cfg, httpClient)
	}
	return
}

func (s *ProviderFactory) getProviderSLS(cfg *config.ProviderConfig, httpClient *http.Client) (cli *sls.Client, err error) {
	cli = &sls.Client{
		Endpoint:   cfg.AuthConfig.AuthForSls.Endpoint,
		HTTPClient: httpClient,
	}
	cli.WithCredentialsProvider(
		sls.NewStaticCredentialsProvider(
			cfg.AuthConfig.AuthForSls.AccessKeyID,
			cfg.AuthConfig.AuthForSls.AccessKeySecret,
			"",
		),
	)
	return
}

func (s *ProviderFactory) getProviderElastic(cfg *config.ProviderConfig, httpClient *http.Client) (cli *elastic.Client, err error) {
	cli, err = elastic.NewClient(
		elastic.SetURL(cfg.AuthConfig.AuthForElasticsearch.URL),
		elastic.SetBasicAuth(cfg.AuthConfig.AuthForElasticsearch.Username, cfg.AuthConfig.AuthForElasticsearch.Password),
		elastic.SetHttpClient(httpClient),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	return
}
