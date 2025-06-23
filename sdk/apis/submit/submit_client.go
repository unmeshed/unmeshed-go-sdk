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
	httpClient         *http.Client
	httpRequestFactory *apis.HttpRequestFactory
	clientConfig       *configs.ClientConfig
	mainQueue          *common.Queue
	retryQueue         *common.Queue
	submitTracker      map[int64]*common.WorkResponseTracker
	submitTrackerLock  sync.Mutex
	stopPolling        atomic.Bool
	workerWg           sync.WaitGroup
	cleanupWg          sync.WaitGroup
}

func NewSubmitClient(httpRequestFactory *apis.HttpRequestFactory, clientConfig *configs.ClientConfig) *SubmitClient {
	if clientConfig.GetClientID() == "" {
		log.Fatal("Cannot submit results without a clientId")
	}

	client := &SubmitClient{
		httpClient:         &http.Client{},
		httpRequestFactory: httpRequestFactory,
		clientConfig:       clientConfig,
		mainQueue:          common.NewQueue(100000),
		retryQueue:         common.NewQueue(100000),
		submitTracker:      make(map[int64]*common.WorkResponseTracker),
	}

	disabled := strings.ToLower(os.Getenv("DISABLE_SUBMIT_CLIENT")) == "true"
	if !disabled {
		for i := 0; i < 5; i++ {
			client.workerWg.Add(1)
			go client.processQueue(client.mainQueue, "main")
			client.workerWg.Add(1)
			go client.processQueue(client.retryQueue, "retry")
		}
		client.cleanupWg.Add(1)
		go client.cleanupLingeringSubmitTrackers()
	}
	return client
}

func (c *SubmitClient) Stop() {
	c.stopPolling.Store(true)
	c.workerWg.Wait()
	c.cleanupWg.Wait()
}

func (c *SubmitClient) cleanupLingeringSubmitTrackers() {
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		c.cleanupWg.Done()
	}()
	for !c.stopPolling.Load() {
		currentMillis := time.Now().UnixMilli()
		c.submitTrackerLock.Lock()
		for stepID, tracker := range c.submitTracker {
			if currentMillis-tracker.QueuedTime > 10*60*1000 {
				tracker.StepPollState.Release(1)
				delete(c.submitTracker, stepID)
			}
		}
		c.submitTrackerLock.Unlock()
		time.Sleep(3 * time.Second)
	}
}

func (c *SubmitClient) processQueue(queue *common.Queue, queueType string) {
	defer c.workerWg.Done()
	timeout := int(c.clientConfig.GetSubmitClientPollTimeoutSeconds())
	if timeout <= 0 {
		timeout = 30
	}
	for !c.stopPolling.Load() {
		var batch []*common.WorkResponse
		itemReceived := false
		start := time.Now()
		for !itemReceived && !c.stopPolling.Load() {
			if !queue.Empty() {
				if workResponse, fetched := queue.Get(); fetched {
					batch = append(batch, workResponse)
					itemReceived = true
					break
				}
			}
			if time.Since(start) > time.Duration(timeout)*time.Second {
				if !c.stopPolling.Load() {
					log.Printf("No item received from queue %s in %d seconds, retrying...", queueType, timeout)
				}
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		for i := 1; i < c.clientConfig.GetResponseSubmitBatchSize(); i++ {
			if !queue.Empty() {
				if workResponse, fetched := queue.Get(); fetched {
					batch = append(batch, workResponse)
				}
			} else {
				break
			}
		}
		if len(batch) == 0 {
			continue
		}
		if err := c.processBatch(batch); err != nil {
			log.Printf("Bulk request failed for batch. Re-queuing all items. Error: %v", err)
			time.Sleep(3 * time.Second)
			for _, workResponse := range batch {
				c.handleAllRequestFailure(workResponse, err.Error())
			}
		}
	}
}

func (c *SubmitClient) processBatch(batch []*common.WorkResponse) error {
	bodyBytes, err := json.Marshal(batch)
	if err != nil {
		return err
	}
	params := map[string]interface{}{}
	resp, err := c.httpRequestFactory.CreatePostRequest(CLIENTS_RESULTS_URL, params, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("response status %d: %s", resp.StatusCode, string(errorBody))
	}
	var responseMap map[string]*common.ClientSubmitResult
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return err
	}
	c.processBatchResults(batch, responseMap)
	return nil
}

func (c *SubmitClient) processBatchResults(batch []*common.WorkResponse, responseMap map[string]*common.ClientSubmitResult) {
	for _, workResponse := range batch {
		stepId := fmt.Sprintf("%d", workResponse.GetStepID())
		c.submitTrackerLock.Lock()
		workResponseTracker := c.submitTracker[workResponse.GetStepID()]
		c.submitTrackerLock.Unlock()
		result, exists := responseMap[stepId]
		if !exists || (result != nil && len(result.GetErrorMessage()) != 0) {
			errorMessage := "No result"
			if result != nil && result.GetErrorMessage() != "" {
				errorMessage = result.GetErrorMessage()
			}
			log.Printf("Error for WorkResponse %d %d: %s", workResponse.GetProcessID(), workResponse.GetStepID(), errorMessage)
			c.enqueueForRetry(workResponse, result, workResponseTracker)
		} else {
			log.Printf("Result from stepId %d %d submitted!", workResponse.GetProcessID(), workResponse.GetStepID())
			c.submitTrackerLock.Lock()
			delete(c.submitTracker, workResponse.GetStepID())
			c.submitTrackerLock.Unlock()
			if workResponseTracker != nil {
				workResponseTracker.StepPollState.Release(1)
			}
		}
	}
}

func (c *SubmitClient) handleAllRequestFailure(workResponse *common.WorkResponse, message string) {
	c.submitTrackerLock.Lock()
	workResponseTracker := c.submitTracker[workResponse.GetStepID()]
	c.submitTrackerLock.Unlock()
	if workResponseTracker != nil {
		c.enqueueForRetry(workResponse, common.NewClientSubmitResult(workResponse.ProcessID, workResponse.StepID, 400, message), workResponseTracker)
	}
}

func (c *SubmitClient) enqueueForRetry(workResponse *common.WorkResponse, result *common.ClientSubmitResult, workResponseTracker *common.WorkResponseTracker) {
	if workResponseTracker == nil {
		return
	}
	if c.isPermanentError(result) {
		log.Printf("Permanent error for WorkResponse %d: %s", workResponse.GetProcessID(), result.GetErrorMessage())
		c.submitTrackerLock.Lock()
		delete(c.submitTracker, workResponse.GetStepID())
		c.submitTrackerLock.Unlock()
		workResponseTracker.StepPollState.Release(1)
		return
	}
	count := workResponseTracker.RetryCount + 1
	if count > int(c.clientConfig.GetMaxSubmitAttempts()) {
		log.Printf("Max retry attempts reached for WorkResponse %d - %d", workResponse.GetStepID(), workResponse.GetProcessID())
		c.submitTrackerLock.Lock()
		delete(c.submitTracker, workResponse.GetStepID())
		c.submitTrackerLock.Unlock()
		workResponseTracker.StepPollState.Release(1)
		return
	}
	workResponseTracker.RetryCount = count
	c.retryQueue.Put(workResponse)
	log.Printf("Re-queued WorkResponse %d for retry attempt %d", workResponse.GetProcessID(), count)
}

func (c *SubmitClient) isPermanentError(result *common.ClientSubmitResult) bool {
	if result == nil || result.GetErrorMessage() == "" {
		return false
	}
	for _, keyword := range c.clientConfig.PermanentErrorKeywords() {
		if strings.Contains(result.GetErrorMessage(), keyword) {
			return true
		}
	}
	return false
}

func (c *SubmitClient) Submit(workResponse *common.WorkResponse, stepPollState *common.StepPollState) error {
	log.Printf("Submitting results to queue: %+v", workResponse)
	c.submitTrackerLock.Lock()
	epochMillis := time.Now().UnixMilli()
	tracker := common.NewWorkResponseTracker(workResponse)
	tracker.QueuedTime = epochMillis
	tracker.StepPollState = stepPollState
	tracker.RetryCount = 0
	c.submitTracker[workResponse.GetStepID()] = tracker
	c.submitTrackerLock.Unlock()
	c.mainQueue.Put(workResponse)
	log.Printf("Result[%v] from stepId %d queued!", workResponse.GetStatus(), workResponse.GetStepID())
	return nil
}

func (c *SubmitClient) GetSubmitTrackerSize() int {
	c.submitTrackerLock.Lock()
	defer c.submitTrackerLock.Unlock()
	return len(c.submitTracker)
}
