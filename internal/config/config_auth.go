package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AuthConfig struct {
	AuthForSls           *AuthForSls           `json:"auth_for_sls,omitempty"`
	AuthForSlsProxy      *AuthForSlsProxy      `json:"auth_for_sls_proxy,omitempty"`
	AuthForElasticsearch *AuthForElasticsearch `json:"auth_for_elasticsearch,omitempty"`
	AuthForKibanaProxy   *AuthForKibanaProxy   `json:"auth_for_kibana_proxy,omitempty"`
}

type AuthForSls struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type AuthForSlsProxy struct {
}

type AuthForElasticsearch struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthForKibanaProxy struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	KibanaAddr string `json:"kibana_addr"`
	KbnVersion string `json:"kbn_version"`
	Cookie     string `json:"cookie"`
}

func (s *AuthForKibanaProxy) init() (err error) {
	type kibanaAuth struct {
		ProviderType string `json:"providerType"`
		ProviderName string `json:"providerName"`
		CurrentURL   string `json:"currentURL"`
		Params       struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"params"`
	}
	auth := kibanaAuth{
		ProviderType: "basic",
		ProviderName: "basic",
		CurrentURL:   s.KibanaAddr + "/login?next=%2Fapp%2Fdiscover#/",
	}
	auth.Params.Username = s.Username
	auth.Params.Password = s.Password
	data, err := json.Marshal(&auth)
	if err != nil {
		return
	}
	request, err := http.NewRequest(http.MethodPost, s.KibanaAddr+"/internal/security/login", bytes.NewReader(data))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("kbn-version", s.KbnVersion)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer func() { _ = response.Body.Close() }()
	var cookies []string
	for _, cookie := range response.Cookies() {
		cookies = append(cookies, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	s.Cookie = strings.Join(cookies, "; ")
	return
}
