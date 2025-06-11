package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/http"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

type RegistrationClient struct {
	clientConfig       *configs.ClientConfig
	httpClient         *http.Client
	requestFactory     *apis.HttpRequestFactory
	workers            []workers.Worker
	clientsRegisterURL string
}

const (
	CLIENTS_REGISTER_URL = "api/clients/register"
)

func NewRegistrationClient(clientConfig *configs.ClientConfig, httpClientFactory *apis.HttpClientFactory,
	httpRequestFactory *apis.HttpRequestFactory) *RegistrationClient {
	return &RegistrationClient{
		clientConfig:       clientConfig,
		httpClient:         httpClientFactory.Create(),
		requestFactory:     httpRequestFactory,
		workers:            []workers.Worker{},
		clientsRegisterURL: CLIENTS_REGISTER_URL,
	}
}

func (rc *RegistrationClient) AddWorkers(workers []workers.Worker) {
	rc.workers = append(rc.workers, workers...)
}

func (rc *RegistrationClient) GetWorkers() []workers.Worker {
	return rc.workers
}

func (rc *RegistrationClient) RenewRegistration() (string, error) {
	supportedSteps := make([]map[string]interface{}, 0)

	for _, worker := range rc.workers {
		step := map[string]interface{}{
			"orgId":     0,
			"namespace": worker.GetNamespace(),
			"stepType":  "WORKER",
			"name":      worker.GetName(),
		}
		supportedSteps = append(supportedSteps, step)
	}

	log.Printf("Renewing registration for the following workers: %v", supportedSteps)

	data, err := json.Marshal(supportedSteps)
	if err != nil {
		return "", fmt.Errorf("failed to marshal supported steps: %w", err)
	}

	params := map[string]interface{}{}
	delay := 2 * time.Second

	for {
		response, err := rc.requestFactory.CreatePutRequest(rc.clientsRegisterURL, params, data)
		if err != nil {
			log.Printf("An error occurred while renewing registration: %v", err)
			log.Printf("Retrying in %d seconds...", int(delay.Seconds()))
			time.Sleep(delay)
			continue
		}

		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read response body: %w", err)
			}
			return string(body), nil
		}

		if response.StatusCode != http.StatusOK {
			errorBody, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("Did not receive 200! Status: %d, Failed to read error response: %v", response.StatusCode, err)
				return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}
			if len(errorBody) > 0 {
				log.Printf("Did not receive 200! Status: %d, Error: %s", response.StatusCode, string(errorBody))
			} else {
				log.Printf("Did not receive 200! Status: %d", response.StatusCode)
			}
			return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
	}
}

func (rc *RegistrationClient) GetWorkerStepNames() []common.WorkerStepName {
	stepNames := make([]common.WorkerStepName, 0)
	for _, worker := range rc.workers {
		stepName := common.WorkerStepName{
			StepQueueNameData: common.StepQueueNameData{
				OrgId:     0,
				Namespace: worker.GetNamespace(),
				StepType:  "WORKER",
				Name:      worker.GetName(),
			},
		}
		stepNames = append(stepNames, stepName)
	}
	return stepNames
}
