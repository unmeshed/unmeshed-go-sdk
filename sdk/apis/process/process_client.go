package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/http"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

type ProcessClient struct {
	httpClientFactory    *apis.HttpClientFactory
	httpRequestFactory   *apis.HttpRequestFactory
	clientConfig         *configs.ClientConfig
	runProcessRequestURL string
}

func NewProcessClient(httpClientFactory *apis.HttpClientFactory, httpRequestFactory *apis.HttpRequestFactory, clientConfig *configs.ClientConfig) *ProcessClient {
	return &ProcessClient{
		httpClientFactory:    httpClientFactory,
		httpRequestFactory:   httpRequestFactory,
		clientConfig:         clientConfig,
		runProcessRequestURL: "api/process/",
	}
}

func (pc *ProcessClient) RunProcessAsync(processRequestData *common.ProcessRequestData) (*common.ProcessData, error) {
	params := map[string]interface{}{
		"clientId": pc.clientConfig.GetClientID(),
	}

	jsonBody, err := json.Marshal(processRequestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal process request data: %w", err)
	}

	response, err := pc.httpRequestFactory.CreatePostRequest(pc.runProcessRequestURL+"runAsync", params, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("invalid process run request: %s", string(body))
	}

	var processData common.ProcessData
	if err := json.NewDecoder(response.Body).Decode(&processData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &processData, nil
}

func (pc *ProcessClient) RunProcessSync(processRequestData *common.ProcessRequestData, processTimeoutSeconds int) (*common.ProcessData, error) {
	params := map[string]interface{}{
		"clientId": pc.clientConfig.GetClientID(),
	}
	if processTimeoutSeconds > 0 {
		params["timeout"] = processTimeoutSeconds
	}

	jsonBody, err := json.Marshal(processRequestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal process request data: %w", err)
	}

	response, err := pc.httpRequestFactory.CreatePostRequest(pc.runProcessRequestURL+"runSync", params, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("invalid process run request: %s", string(body))
	}

	var processData common.ProcessData
	if err := json.NewDecoder(response.Body).Decode(&processData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &processData, nil
}

func (pc *ProcessClient) GetProcessData(processID int64, includeSteps bool) (*common.ProcessData, error) {
	if processID == 0 {
		return nil, fmt.Errorf("process ID cannot be zero")
	}

	url := fmt.Sprintf("%scontext/%d", pc.runProcessRequestURL, processID)
	params := map[string]interface{}{
		"includeSteps": includeSteps,
	}

	response, err := pc.httpRequestFactory.CreateGetRequest(url, params)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("invalid fetch process data request: %s", string(body))
	}

	var processData common.ProcessData
	if err := json.NewDecoder(response.Body).Decode(&processData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &processData, nil
}

func (pc *ProcessClient) GetStepData(stepID int64) (*common.StepData, error) {
	if stepID == 0 {
		return nil, fmt.Errorf("step ID cannot be zero")
	}

	url := fmt.Sprintf("%sstepContext/%d", pc.runProcessRequestURL, stepID)
	response, err := pc.httpRequestFactory.CreateGetRequest(url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("invalid fetch step data request: %s", string(body))
	}

	var stepData common.StepData
	if err := json.NewDecoder(response.Body).Decode(&stepData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &stepData, nil
}

func (pc *ProcessClient) BulkTerminate(processIDs []int64, reason string) (*common.ProcessActionResponseData, error) {
	if len(processIDs) == 0 {
		return nil, fmt.Errorf("process IDs cannot be empty")
	}

	url := pc.runProcessRequestURL + "bulkTerminate"
	params := map[string]interface{}{}
	if reason != "" {
		params["reason"] = reason
	}

	jsonBody, err := json.Marshal(processIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal process IDs: %w", err)
	}

	response, err := pc.httpRequestFactory.CreatePostRequest(url, params, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to bulk terminate: %s", string(body))
	}

	var responseData common.ProcessActionResponseData
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &responseData, nil
}

func (pc *ProcessClient) BulkResume(processIDs []int64) (*common.ProcessActionResponseData, error) {
	if len(processIDs) == 0 {
		return nil, fmt.Errorf("process IDs cannot be empty")
	}

	url := pc.runProcessRequestURL + "bulkResume"
	params := map[string]interface{}{
		"clientId": pc.clientConfig.GetClientID(),
	}

	jsonBody, err := json.Marshal(processIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal process IDs: %w", err)
	}

	response, err := pc.httpRequestFactory.CreatePostRequest(url, params, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to bulk resume: %s", string(body))
	}

	var responseData common.ProcessActionResponseData
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &responseData, nil
}

func (pc *ProcessClient) BulkReviewed(processIDs []int64, reason string) (*common.ProcessActionResponseData, error) {
	if len(processIDs) == 0 {
		return nil, fmt.Errorf("process IDs cannot be empty")
	}

	url := pc.runProcessRequestURL + "bulkReviewed"
	params := map[string]interface{}{
		"clientId": pc.clientConfig.GetClientID(),
	}
	if reason != "" {
		params["reason"] = reason
	}

	jsonBody, err := json.Marshal(processIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal process IDs: %w", err)
	}

	response, err := pc.httpRequestFactory.CreatePostRequest(url, params, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to bulk reviewed: %s", string(body))
	}

	var responseData common.ProcessActionResponseData
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &responseData, nil
}

func (pc *ProcessClient) Rerun(processID int64, version int) (*common.ProcessData, error) {
	if processID == 0 {
		return nil, fmt.Errorf("process ID cannot be zero")
	}

	params := map[string]interface{}{
		"clientId":  pc.clientConfig.GetClientID(),
		"processId": processID,
	}
	if version > 0 {
		params["version"] = version
	}

	url := pc.runProcessRequestURL + "rerun"
	response, err := pc.httpRequestFactory.CreatePostRequest(url, params, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to rerun request: %s", string(body))
	}

	var processData common.ProcessData
	if err := json.NewDecoder(response.Body).Decode(&processData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &processData, nil
}

func (pc *ProcessClient) SearchProcessExecutions(params *common.ProcessSearchRequest) ([]*common.ProcessData, error) {
	queryParams := make(map[string]interface{})

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search params: %w", err)
	}

	if err := json.Unmarshal(jsonData, &queryParams); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search params: %w", err)
	}

	filteredParams := make(map[string]interface{})
	for k, v := range queryParams {
		if v != nil {
			switch val := v.(type) {
			case []interface{}:
				if len(val) > 0 {
					strValues := make([]string, len(val))
					for i, item := range val {
						strValues[i] = fmt.Sprintf("%v", item)
					}
					filteredParams[k] = strings.Join(strValues, ",")
				}
			default:
				filteredParams[k] = v
			}
		}
	}

	url := "api/stats/process/search"
	response, err := pc.httpRequestFactory.CreateGetRequest(url, filteredParams)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("invalid fetch processes data: %s", string(body))
	}

	var processesData []*common.ProcessData
	if err := json.NewDecoder(response.Body).Decode(&processesData); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return processesData, nil
}

func (pc *ProcessClient) InvokeAPIMappingGet(endpoint string, id string, correlationID string, apiCallType common.ApiCallType) (map[string]interface{}, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	queryParams := make(map[string]interface{})
	if id != "" {
		queryParams["id"] = id
	}
	if correlationID != "" {
		queryParams["correlationId"] = correlationID
	}
	if apiCallType != "" {
		queryParams["apiCallType"] = apiCallType.String()
	} else {
		queryParams["apiCallType"] = common.ApiCallTypeAsync.String()
	}

	url := "api/call/" + endpoint
	response, err := pc.httpRequestFactory.CreateGetRequest(url, queryParams)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed invoking webhook get request: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return result, nil
}

func (pc *ProcessClient) InvokeAPIMappingPost(endpoint string, input map[string]interface{}, id string, correlationID string, apiCallType common.ApiCallType) (map[string]interface{}, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	queryParams := make(map[string]interface{})
	if id != "" {
		queryParams["id"] = id
	}
	if correlationID != "" {
		queryParams["correlationId"] = correlationID
	}
	if apiCallType != "" {
		queryParams["apiCallType"] = apiCallType.String()
	} else {
		queryParams["apiCallType"] = common.ApiCallTypeAsync.String()
	}

	jsonBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	url := "api/call/" + endpoint
	response, err := pc.httpRequestFactory.CreatePostRequest(url, queryParams, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed invoking webhook post request: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return result, nil
}
