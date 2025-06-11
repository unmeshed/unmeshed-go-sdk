package apis

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/utils"
)

type HttpRequestFactory struct {
	clientConfig *configs.ClientConfig
	baseURL      string
	port         int
	bearerValue  string
	client       *utils.RetryClient
}

func NewHttpRequestFactory(clientConfig *configs.ClientConfig) *HttpRequestFactory {
	hashedToken, err := utils.CreateSecureHash(clientConfig.AuthToken)
	if err != nil {
		fmt.Errorf("error creating secure hash for token: %w", err)
	}

	bearerValue := fmt.Sprintf("Bearer client.sdk.%s.%s",
		clientConfig.ClientID,
		hashedToken,
	)

	// Configure transport with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        10,               // Maximum number of idle connections across all hosts (pool_maxsize)
		MaxIdleConnsPerHost: 2,                // Maximum number of idle connections per host (pool_connections)
		MaxConnsPerHost:     10,               // Maximum number of connections per host (pool_maxsize)
		IdleConnTimeout:     90 * time.Second, // How long an idle connection is kept in the pool
		DisableCompression:  false,            // Enable compression
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Set a reasonable timeout
	}

	// Add retry logic
	retryClient := &utils.RetryClient{
		Client:     client,
		MaxRetries: 3, // max_retries
		RetryDelay: time.Second,
	}

	return &HttpRequestFactory{
		clientConfig: clientConfig,
		baseURL:      clientConfig.GetBaseURL(),
		port:         clientConfig.GetPort(),
		bearerValue:  bearerValue,
		client:       retryClient,
	}
}

func (factory *HttpRequestFactory) buildURI(path string, params map[string]interface{}) string {
	baseURL := factory.baseURL

	baseURL = strings.TrimSuffix(baseURL, "/")

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		panic("Invalid base URL")
	}

	hasPort := parsedURL.Port() != ""
	var urlStr string
	if hasPort || strings.HasPrefix(baseURL, "https:") {
		urlStr = fmt.Sprintf("%s/%s", baseURL, path)
	} else {
		urlStr = fmt.Sprintf("%s:%d/%s", baseURL, factory.port, path)
	}

	if params != nil {
		queryString := url.QueryEscape(fmt.Sprintf("%v", params))
		urlStr = fmt.Sprintf("%s?%s", urlStr, queryString)
	}

	return urlStr
}

func (factory *HttpRequestFactory) CreateGetRequest(path string, params map[string]interface{}) (*http.Response, error) {
	uri := factory.buildURI(path, params)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", factory.bearerValue)
	return factory.client.Do(req)
}

func (factory *HttpRequestFactory) CreatePostRequest(path string, params map[string]interface{}, body []byte) (*http.Response, error) {
	uri := factory.buildURI(path, params)
	requestBody := bytes.NewReader(body)
	req, err := http.NewRequest("POST", uri, requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", factory.bearerValue)
	return factory.client.Do(req)
}

func (factory *HttpRequestFactory) CreatePutRequest(path string, params map[string]interface{}, body []byte) (*http.Response, error) {
	uri := factory.buildURI(path, params)
	requestBody := bytes.NewReader(body)
	req, err := http.NewRequest("PUT", uri, requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", factory.bearerValue)
	return factory.client.Do(req)
}

func (factory *HttpRequestFactory) CreatePostRequestWithBody(path string, body []byte) (*http.Response, error) {
	return factory.CreatePostRequest(path, nil, body)
}
