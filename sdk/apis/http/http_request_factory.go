package apis

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	baseURL := strings.TrimSuffix(factory.baseURL, "/")

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		panic(fmt.Sprintf("Invalid base URL: %v", err))
	}

	// Format base path with port if not already present
	hasPort := parsedURL.Port() != ""
	var urlStr string
	if hasPort || strings.HasPrefix(baseURL, "https:") {
		urlStr = fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(path, "/"))
	} else {
		urlStr = fmt.Sprintf("%s:%d/%s", baseURL, factory.port, strings.TrimPrefix(path, "/"))
	}

	if len(params) > 0 {
		query := url.Values{}
		for key, value := range params {
			if value == nil {
				continue
			}

			switch v := value.(type) {
			case string:
				query.Add(key, v)
			case fmt.Stringer:
				query.Add(key, v.String())
			case int:
				query.Add(key, strconv.Itoa(v))
			case int32:
				query.Add(key, strconv.FormatInt(int64(v), 10))
			case int64:
				query.Add(key, strconv.FormatInt(v, 10))
			case float32:
				query.Add(key, strconv.FormatFloat(float64(v), 'f', -1, 32))
			case float64:
				query.Add(key, strconv.FormatFloat(v, 'f', -1, 64))
			case bool:
				query.Add(key, strconv.FormatBool(v))
			case []string:
				for _, s := range v {
					query.Add(key, s)
				}
			case []interface{}:
				for _, item := range v {
					if item == nil {
						continue
					}
					query.Add(key, fmt.Sprintf("%v", item))
				}
			default:
				query.Add(key, fmt.Sprintf("%v", value))
			}
		}
		encoded := query.Encode()
		if encoded != "" {
			urlStr = fmt.Sprintf("%s?%s", urlStr, encoded)
		}
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

func (factory *HttpRequestFactory) CreatePostRequestWithHeaders(path string, params map[string]interface{}, headers map[string]string, body []byte) (*http.Response, error) {
	uri := factory.buildURI(path, params)
	requestBody := bytes.NewReader(body)
	req, err := http.NewRequest("POST", uri, requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", factory.bearerValue)

	for k, v := range headers {
       req.Header.Set(k, v)
    }
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

func (factory *HttpRequestFactory) CreateDeleteRequest(path string, params map[string]interface{}, body []byte) (*http.Response, error) {
	uri := factory.buildURI(path, params)
	var requestBody *bytes.Reader

	if body != nil {
		requestBody = bytes.NewReader(body)
	} else {
		requestBody = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest("DELETE", uri, requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", factory.bearerValue)

	return factory.client.Do(req)
}
