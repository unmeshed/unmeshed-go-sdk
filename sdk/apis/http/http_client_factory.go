package apis

import (
	"net/http"
	"time"

	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

type HttpClientFactory struct {
	clientConfig *configs.ClientConfig
}

func NewHttpClientFactory(clientConfig *configs.ClientConfig) *HttpClientFactory {
	return &HttpClientFactory{clientConfig: clientConfig}
}

func (factory *HttpClientFactory) Create() *http.Client {
	timeoutSecs := factory.clientConfig.ConnectionTimeoutSecs

	timeout := time.Duration(timeoutSecs) * time.Second
	client := &http.Client{
		Timeout: timeout,
	}
	return client
}
