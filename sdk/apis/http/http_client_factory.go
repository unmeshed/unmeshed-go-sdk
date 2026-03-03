package apis

import (
	"crypto/tls"
	"net/http"
	"os"
	"strings"
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
	transport := &http.Transport{}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("DISABLE_SSL_VERIFICATION")), "true") {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	return client
}
