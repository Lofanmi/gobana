package sdk_sls

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
)

type Client struct {
	Client sls.ClientInterface
}

// NewClient
// @autowire(set=sdks)
func NewClient() *Client {
	Endpoint := ""
	AccessKeyId := ""
	AccessKeySecret := ""
	SecurityToken := ""
	_ = sls.CreateNormalInterface(Endpoint, AccessKeyId, AccessKeySecret, SecurityToken)
	return nil
}
