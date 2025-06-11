package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/http"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

type PollerClient struct {
	clientConfig       *configs.ClientConfig
	httpClient         *http.Client
	httpRequestFactory *apis.HttpRequestFactory
	CLIENTS_POLL_URL   string
}

func NewPollerClient(clientConfig *configs.ClientConfig, httpClientFactory *apis.HttpClientFactory, httpRequestFactory *apis.HttpRequestFactory) *PollerClient {
	clientPollURL := "api/clients/poll"
	return &PollerClient{
		clientConfig:       clientConfig,
		httpClient:         httpClientFactory.Create(),
		httpRequestFactory: httpRequestFactory,
		CLIENTS_POLL_URL:   clientPollURL,
	}
}

func (pc *PollerClient) Poll(stepSizes []common.StepSize) ([]common.WorkRequest, error) {
	if len(stepSizes) == 0 {
		stepSizes = []common.StepSize{}
	}

	jsonData, err := json.Marshal(stepSizes)
	if err != nil {
		log.Printf("Error marshalling JSON body: %v", err)
		return nil, fmt.Errorf("failed to marshal JSON body: %w", err)
	}

	params := map[string]interface{}{
		"size": fmt.Sprintf("%d", pc.clientConfig.GetWorkRequestBatchSize()),
	}

	clientPollUrl := pc.CLIENTS_POLL_URL

	response, err := pc.httpRequestFactory.CreatePostRequest(clientPollUrl, params, jsonData)
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Did not receive 200! Status: %d, Failed to read error response: %v", response.StatusCode, err)
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
		if len(errorBody) > 0 {
			log.Printf("Did not receive 200! Status: %d, Error: %s", response.StatusCode, string(errorBody))
		} else {
			log.Printf("Did not receive 200! Status: %d", response.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var workRequests []common.WorkRequest
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&workRequests); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return nil, fmt.Errorf("failed to decode response JSON: %w", err)
	}

	return workRequests, nil
}
