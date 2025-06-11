package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/http"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

const (
	CLIENTS_RESULTS_URL = "api/clients/bulkResults"
	MAX_RETRIES         = 3
	INITIAL_BACKOFF     = 100 * time.Millisecond
	MAX_BACKOFF         = 5 * time.Second
)

type SubmitClient struct {
	HttpClient         *http.Client
	HttpRequestFactory *apis.HttpRequestFactory
	ClientConfig       *configs.ClientConfig
	mainQueue          *common.Queue
	retryQueue         *common.Queue
	submitTracker      map[int64]*common.WorkResponseTracker
	submitTrackerLock  sync.RWMutex
	workerWg           sync.WaitGroup
	stopChan           chan struct{}
	disabled           atomic.Bool
	activeWorkers      atomic.Int32
}

func NewSubmitClient(httpRequestFactory *apis.HttpRequestFactory, clientConfig *configs.ClientConfig) *SubmitClient {
	// Check if submit client is disabled via environment variable
	disabled := strings.ToLower(os.Getenv("DISABLE_SUBMIT_CLIENT")) == "true"

	client := &SubmitClient{
		HttpClient:         &http.Client{},
		HttpRequestFactory: httpRequestFactory,
		ClientConfig:       clientConfig,
		mainQueue:          common.NewQueue(10000),
		retryQueue:         common.NewQueue(10000),
		submitTracker:      make(map[int64]*common.WorkResponseTracker),
		submitTrackerLock:  sync.RWMutex{},
		stopChan:           make(chan struct{}),
		disabled:           atomic.Bool{},
		activeWorkers:      atomic.Int32{},
	}

	client.disabled.Store(disabled)

	// Validate client ID
	if clientConfig.GetClientID() == "" {
		log.Fatal("Cannot submit results without a clientId")
	}

	if !disabled {
		workersCount := int(3)
		if workersCount < 10 {
			workersCount = 10
		}

		mainWorkers := (workersCount * 2) / 3
		for i := 0; i < mainWorkers; i++ {
			client.workerWg.Add(1)
			client.activeWorkers.Add(1)
			go client.processQueue(client.mainQueue)
		}

		// Start worker for retry queue (1/3 of threads)
		retryWorkers := workersCount - mainWorkers
		for i := 0; i < retryWorkers; i++ {
			client.workerWg.Add(1)
			client.activeWorkers.Add(1)
			go client.processQueue(client.retryQueue)
		}

		go client.cleanupLingeringSubmitTrackers()
	}

	return client
}

func (c *SubmitClient) processQueue(queue *common.Queue) {
	defer func() {
		c.workerWg.Done()
		c.activeWorkers.Add(-1)
	}()

	backoff := INITIAL_BACKOFF
	for {
		select {
		case <-c.stopChan:
			return
		default:
			var batch []*common.WorkResponse
			batchSize := c.ClientConfig.GetResponseSubmitBatchSize()

			for i := 0; i < batchSize; i++ {
				if !queue.Empty() {
					if workResponse, fetched := queue.Get(); fetched {
						batch = append(batch, workResponse)
					}
				}
			}

			if len(batch) == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if err := c.processBatch(batch); err != nil {
				backoff = min(backoff*2, MAX_BACKOFF)
				time.Sleep(backoff)
				continue
			}

			backoff = INITIAL_BACKOFF
		}
	}
}

func (c *SubmitClient) processBatch(batch []*common.WorkResponse) error {
	bodyBytes, err := json.Marshal(batch)
	if err != nil {
		log.Printf("Failed to marshal request body: %v", err)
		return err
	}

	params := map[string]interface{}{}
	resp, err := c.HttpRequestFactory.CreatePostRequest(CLIENTS_RESULTS_URL, params, bodyBytes)
	if err != nil {
		log.Printf("Bulk request failed for batch. Error: %v", err)
		c.handleBatchFailure(batch, err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Bulk request failed with status %d, Failed to read error response: %v", resp.StatusCode, err)
			c.handleBatchFailure(batch, fmt.Sprintf("Response status %d", resp.StatusCode))
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		if len(errorBody) > 0 {
			log.Printf("Bulk request failed with status %d. Error response: %s", resp.StatusCode, string(errorBody))
			c.handleBatchFailure(batch, fmt.Sprintf("Response status %d: %s", resp.StatusCode, string(errorBody)))
		} else {
			log.Printf("Bulk request failed with status %d", resp.StatusCode)
			c.handleBatchFailure(batch, fmt.Sprintf("Response status %d", resp.StatusCode))
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var responseMap map[string]*common.ClientSubmitResult
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		log.Printf("Failed to decode response body: %v", err)
		return err
	}

	c.processBatchResults(batch, responseMap)
	return nil
}

func (c *SubmitClient) processBatchResults(batch []*common.WorkResponse, responseMap map[string]*common.ClientSubmitResult) {
	c.submitTrackerLock.Lock()
	defer c.submitTrackerLock.Unlock()

	for _, workResponse := range batch {
		stepId := fmt.Sprintf("%d", workResponse.GetStepID())
		result, exists := responseMap[stepId]
		tracker := c.submitTracker[workResponse.GetStepID()]

		if !exists || len(result.GetErrorMessage()) != 0 {
			if tracker != nil && tracker.RetryCount < MAX_RETRIES {
				tracker.RetryCount++
				c.retryQueue.Put(workResponse)
			} else {
				c.handlePermanentFailure(workResponse, result, tracker)
			}
		} else {
			c.handleSuccess(workResponse, tracker)
		}
	}
}

func (c *SubmitClient) handleBatchFailure(batch []*common.WorkResponse, message string) {
	c.submitTrackerLock.Lock()
	defer c.submitTrackerLock.Unlock()

	for _, workResponse := range batch {
		if tracker, exists := c.submitTracker[workResponse.GetStepID()]; exists {
			if tracker.RetryCount < MAX_RETRIES {
				tracker.RetryCount++
				c.retryQueue.Put(workResponse)
			} else {
				c.handlePermanentFailure(workResponse, common.NewClientSubmitResult(workResponse.ProcessID, workResponse.StepID, 400, message), tracker)
			}
		}
	}
}

func (c *SubmitClient) handleSuccess(workResponse *common.WorkResponse, tracker *common.WorkResponseTracker) {
	delete(c.submitTracker, workResponse.GetStepID())
	if tracker != nil {
		tracker.StepPollState.Release(1)
	}
}

func (c *SubmitClient) handlePermanentFailure(workResponse *common.WorkResponse, result *common.ClientSubmitResult, tracker *common.WorkResponseTracker) {
	log.Printf("[ERROR:] Permanent error for WorkResponse %s:%s", workResponse.GetProcessID(), result.GetErrorMessage())
	delete(c.submitTracker, workResponse.GetStepID())
	if tracker != nil {
		tracker.StepPollState.Release(1)
	}
}

func (c *SubmitClient) cleanupLingeringSubmitTrackers() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			currentMillis := time.Now().UnixMilli()
			c.submitTrackerLock.Lock()
			for stepID, tracker := range c.submitTracker {
				if currentMillis-tracker.QueuedTime > 10*60*1000 {
					tracker.StepPollState.Release(1)
					delete(c.submitTracker, stepID)
				}
			}
			c.submitTrackerLock.Unlock()
		}
	}
}

func (c *SubmitClient) Submit(workResponse *common.WorkResponse, stepPollState *common.StepPollState) error {
	if c.disabled.Load() {
		stepPollState.Release(stepPollState.GetTotalCount())
		return nil
	}

	c.submitTrackerLock.Lock()
	defer c.submitTrackerLock.Unlock()

	tracker := common.NewWorkResponseTracker(workResponse)
	tracker.QueuedTime = time.Now().UnixMilli()
	tracker.StepPollState = stepPollState
	tracker.RetryCount = 0

	c.submitTracker[workResponse.GetStepID()] = tracker
	c.mainQueue.Put(workResponse)
	return nil
}

func (c *SubmitClient) GetSubmitTrackerSize() int {
	c.submitTrackerLock.RLock()
	defer c.submitTrackerLock.RUnlock()
	return len(c.submitTracker)
}

func (c *SubmitClient) GetActiveWorkers() int32 {
	return c.activeWorkers.Load()
}

func (c *SubmitClient) Close() {
	close(c.stopChan)
	c.workerWg.Wait()
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
