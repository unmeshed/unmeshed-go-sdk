package apis

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/http"
	poller "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/poller"
	process "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/process"
	register "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/register"
	workerRunner "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/runner"
	submit "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/submit"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	workersApi "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

func setupLogging() {
	enableFileLogging := os.Getenv("ENABLE_FILE_LOGGING") == "true"

	if enableFileLogging {
		logsDir := "logs"
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			log.Printf("Failed to create logs directory: %v", err)
			return
		}

		timestamp := time.Now().Format("2006-01-02_15-04-05")
		logFile := filepath.Join(logsDir, fmt.Sprintf("unmeshed_%s.log", timestamp))

		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open log file: %v", err)
			return
		}

		log.SetOutput(file)
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		log.Printf("Logging initialized. Log file: %s", logFile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	}
}

type UnmeshedClient struct {
	ClientConfig             *configs.ClientConfig
	Workers                  []workersApi.Worker
	pollStates               map[string]*common.StepPollState
	permitsZero              atomic.Int32
	pollSizeZeroTime         atomic.Int64
	executingCount           atomic.Int32
	lastPrinted              atomic.Int32
	pollingErrorRepoted      atomic.Bool
	jitterUpperBoundInMillis int
	initialDelayMillis       int
	backoffMultiplier        int
	retryCount               atomic.Int32
	workerByName             map[string]bool
	workerByNameMutex        sync.RWMutex
	httpClientFactory        *apis.HttpClientFactory
	httpRequestFactory       *apis.HttpRequestFactory
	registrationClient       *register.RegistrationClient
	pollerClient             *poller.PollerClient
	submitClient             *submit.SubmitClient
	processClient            *process.ProcessClient
	workResponseBuilder      *common.WorkResponseBuilder
	workerRunner             *workerRunner.WorkerRunner
	stopPolling              atomic.Bool
	done                     chan struct{}

	lastPrintedPolling int64
	lastPrintedRunning int64

	stopOnce sync.Once
}

func NewUnmeshedClient(
	clientConfig *configs.ClientConfig,
) (*UnmeshedClient, error) {
	// Setup logging first
	setupLogging()

	// Validate client ID and token
	if clientConfig.GetClientID() == "" || !clientConfig.HasToken() {
		return nil, fmt.Errorf("cannot initialize without a valid clientId and token")
	}

	httpClientFactory := apis.NewHttpClientFactory(clientConfig)
	httpRequestFactory := apis.NewHttpRequestFactory(clientConfig)
	pollerClient := poller.NewPollerClient(clientConfig, httpClientFactory, httpRequestFactory)
	submitClient := submit.NewSubmitClient(httpRequestFactory, clientConfig)
	processClient := process.NewProcessClient(httpClientFactory, httpRequestFactory, clientConfig)

	unmeshedClient := &UnmeshedClient{
		ClientConfig:             clientConfig,
		Workers:                  []workersApi.Worker{},
		pollStates:               make(map[string]*common.StepPollState),
		permitsZero:              atomic.Int32{},
		pollSizeZeroTime:         atomic.Int64{},
		executingCount:           atomic.Int32{},
		lastPrinted:              atomic.Int32{},
		pollingErrorRepoted:      atomic.Bool{},
		jitterUpperBoundInMillis: 2000,
		initialDelayMillis:       0,
		backoffMultiplier:        2,
		retryCount:               atomic.Int32{},
		workerByName:             make(map[string]bool),
		workerByNameMutex:        sync.RWMutex{},
		httpClientFactory:        httpClientFactory,
		httpRequestFactory:       httpRequestFactory,
		registrationClient:       register.NewRegistrationClient(clientConfig, httpClientFactory, httpRequestFactory),
		submitClient:             submitClient,
		processClient:            processClient,
		pollerClient:             pollerClient,
		workResponseBuilder:      common.NewWorkResponseBuilder(),
		workerRunner:             workerRunner.NewWorkerRunner(),
		stopPolling:              atomic.Bool{},
		done:                     make(chan struct{}),
		lastPrintedPolling:       0,
		lastPrintedRunning:       0,
	}

	return unmeshedClient, nil
}

func (uc *UnmeshedClient) getWorkers() []workersApi.Worker {
	return uc.Workers
}

func formattedWorkerID(namespace string, name string) string {
	return fmt.Sprintf("%s:-#-:%s", namespace, name)
}

func (uc *UnmeshedClient) pollForWork(disableLogRunningWorkerDetails bool) ([]common.WorkRequest, error) {

	workers := uc.registrationClient.GetWorkers()
	workerTasks := []common.StepSize{}
	workerRequestCount := make(map[string]int)

	for _, worker := range workers {
		stepQueueNameData := common.StepQueueNameData{
			OrgId:     0,
			Namespace: worker.GetNamespace(),
			StepType:  "WORKER",
			Name:      worker.GetName(),
		}
		workerId := formattedWorkerID(worker.GetNamespace(), worker.GetName())
		state, exists := uc.pollStates[workerId]
		if !exists {
			return nil, fmt.Errorf("unexpected missing poll state for worker: %s", workerId)
		}
		size := state.AcquireMaxAvailable()
		workerRequestCount[workerId] = size
		if size > 0 {
			workerTask := common.NewStepSize(stepQueueNameData, size)
			workerTasks = append(workerTasks, workerTask)
		}
	}

	now := time.Now().Unix()
	if now-uc.lastPrintedPolling > 10 {
		log.Printf("Tasks being polled: %v", workerTasks)
		uc.lastPrintedPolling = now
	}

	if len(workerTasks) == 0 {
		return nil, nil
	}

	workRequests, err := uc.pollerClient.Poll(workerTasks)
	if err != nil {
		uc.releaseUnusedPermits(make(map[string]int), workerRequestCount)
		return nil, fmt.Errorf("failed to poll work requests: %w", err)
	}

	if len(workRequests) > 0 {
		log.Printf("Received work requests: %d", len(workRequests))
	}

	workerReceivedCount := make(map[string]int)
	for _, workRequest := range workRequests {
		uc.executingCount.Add(1)
		workerId := formattedWorkerID(workRequest.GetStepNamespace(), workRequest.GetStepName())
		workerReceivedCount[workerId]++
	}
	uc.releaseUnusedPermits(workerReceivedCount, workerRequestCount)

	if now-uc.lastPrintedRunning > 10 {
		logEntries := make([]string, 0, len(workers))
		for _, s := range workers {
			workerId := formattedWorkerID(s.GetNamespace(), s.GetName())
			pollState := uc.pollStates[workerId]
			available := pollState.MaxAvailable()
			total := pollState.GetTotalCount()
			requested := workerRequestCount[workerId]
			if !disableLogRunningWorkerDetails {
				logEntries = append(logEntries,
					fmt.Sprintf("%s:%s = Available[%d] / [%d] / [%d]", s.GetNamespace(), s.GetName(), available, requested, total))
			}
		}
		logStr := strings.Join(logEntries, ", ")
		executingCount := uc.executingCount.Load()
		submitTrackerSize := int32(uc.submitClient.GetSubmitTrackerSize())
		if !disableLogRunningWorkerDetails {
			log.Printf("Running : %d st: %d t: %d - permits %s", executingCount, submitTrackerSize, executingCount+submitTrackerSize, logStr)
		}
		uc.lastPrintedRunning = now
	}

	return workRequests, nil
}

func (uc *UnmeshedClient) runStep(worker *workersApi.Worker, workRequest *common.WorkRequest) {
	result, err := uc.workerRunner.RunWorker(worker, workRequest)

	var stepResult *common.StepResult

	if sr, ok := result.(*common.StepResult); ok {
		stepResult = sr
	} else {
		if result == nil {
			result = map[string]interface{}{}
		}
		stepResult = common.NewStepResult(result)
	}

	if err != nil {
		uc.handleWorkCompletion(workRequest, stepResult, &err)
	} else {
		uc.handleWorkCompletion(workRequest, stepResult, nil)
	}
}

func (uc *UnmeshedClient) handleWorkCompletion(workRequest *common.WorkRequest, stepResult *common.StepResult, throwable *error) {
	stepId := formattedWorkerID(workRequest.GetStepNamespace(), workRequest.GetStepName())
	state := uc.pollStates[stepId]

	var workResponse *common.WorkResponse

	if throwable != nil {
		workResponse = uc.workResponseBuilder.FailResponse(workRequest, *throwable)
	} else if stepResult.KeepRunning && stepResult.RescheduleAfterSeconds > 0 {
		workResponse = uc.workResponseBuilder.RunningResponse(workRequest, stepResult)
	} else {
		workResponse = uc.workResponseBuilder.SuccessResponse(workRequest, stepResult)
	}

	uc.submitClient.Submit(workResponse, state)
	uc.executingCount.Add(-1)
}

func (client *UnmeshedClient) releaseUnusedPermits(workerReceivedCount, workerRequestCount map[string]int) {
	for workerID, requestedCount := range workerRequestCount {
		pollState, exists := client.pollStates[workerID]

		if exists {
			receivedCount := workerReceivedCount[workerID]
			pollState.Release(requestedCount - receivedCount)
		}
	}
}

func (client *UnmeshedClient) startAsyncTaskProcessing() {
	const logInterval = 60 * time.Second
	const minBackoff = 100 * time.Millisecond
	const maxBackoff = 20 * time.Second

	disableLogRunningWorkerDetails := os.Getenv("DISABLE_RUNNING_WORKER_LOGS") == "true"

	// Determine worker pool size
	workerCount := int(client.ClientConfig.GetMaxWorkers())
	if workerCount < 10 {
		workerCount = 10
	}
	workQueue := make(chan common.WorkRequest, workerCount*2)

	// Start worker pool
	var workerWg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for workRequest := range workQueue {
				var foundWorker *workersApi.Worker
				for i := range client.Workers {
					if client.Workers[i].GetName() == workRequest.GetStepName() {
						foundWorker = &client.Workers[i]
						break
					}
				}
				if foundWorker != nil {
					client.runStep(foundWorker, &workRequest)
				} else {
					log.Printf("No worker found for step '%s'\n", workRequest.GetStepName())
				}
			}
		}()
	}

	go func() {
		var (
			lastLogTime    = time.Now()
			pollRetryCount = 1
		)
		for !client.stopPolling.Load() {
			pollInterval := time.Duration(client.ClientConfig.GetDelayMillis()) * time.Millisecond
			workRequests, err := client.pollForWork(disableLogRunningWorkerDetails)

			if err != nil {
				backoff := minBackoff << (pollRetryCount - 1)
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				log.Printf("Polling error: %v, will retry after %v", err, backoff)
				pollRetryCount++
				time.Sleep(backoff)
				continue
			} else {
				pollRetryCount = 1
			}

			if len(workRequests) > 0 {
				for _, workRequest := range workRequests {
					workQueue <- workRequest
				}
				log.Printf("All tasks scheduled. Continuing the polling")
			}

			if time.Since(lastLogTime) >= logInterval {
				log.Printf("Poll interval is %d ms", client.ClientConfig.GetDelayMillis())
				lastLogTime = time.Now()
			}

			time.Sleep(pollInterval)
		}
		close(workQueue)
	}()

	// Wait for all workers to finish before returning (when stopPolling is set)
	go func() {
		workerWg.Wait()
		log.Println("Worker pool exited gracefully.")
	}()
}

func (uc *UnmeshedClient) renewRegistrationWithRetry(renewRegistrationTask interface{}) (string, error) {
	const delay = 2 * time.Second

	for {
		log.Printf("Attempting to renew registration")

		results := reflect.ValueOf(renewRegistrationTask).Call(nil)
		if len(results) != 2 {
			return "", fmt.Errorf("expected 2 return values from renewRegistrationTask")
		}

		responseText := results[0].Interface().(string)
		errInterface := results[1].Interface()

		if errInterface != nil {
			if err, ok := errInterface.(error); ok {
				log.Printf("An error occurred while renewing registration: %v", err)
				log.Printf("Retrying in %d seconds...", int(delay.Seconds()))
				time.Sleep(delay)
				continue
			}
		}

		log.Printf("Successfully renewed registration for workers.")
		return responseText, nil
	}
}

func (uc *UnmeshedClient) Start() {
	if !uc.ClientConfig.HasToken() {
		log.Fatal("Credentials not configured correctly. Client configuration requires auth client id and token to be set.")
	}

	if uc.registrationClient.GetWorkers() == nil || len(uc.registrationClient.GetWorkers()) == 0 {
		log.Printf("No workers configured. Will not poll for any work.")
		return
	}

	for _, worker := range uc.registrationClient.GetWorkers() {
		defaultMaxSize := worker.GetMaxInProgress()
		workerId := formattedWorkerID(worker.GetNamespace(), worker.GetName())

		uc.pollStates[workerId] = common.NewStepPollState(defaultMaxSize)
	}

	log.Printf("Registering %d workers", len(uc.registrationClient.GetWorkers()))

	renewRegistrationTask := uc.registrationClient.RenewRegistration
	_, err := uc.renewRegistrationWithRetry(renewRegistrationTask)
	if err != nil {
		log.Printf("Error renewing registration: %v", err)
	}

	log.Printf("Unmeshed Go SDK started")
	go uc.startAsyncTaskProcessing()
	<-uc.done
}

func (uc *UnmeshedClient) registerWorker(worker *workersApi.Worker) error {
	uc.workerByNameMutex.Lock()
	defer uc.workerByNameMutex.Unlock()

	if _, exists := uc.workerByName[worker.GetName()]; exists {
		return fmt.Errorf("worker with name %s is already registered", worker.GetName())
	}

	method := worker.GetExecutionMethod()
	if method == nil {
		return fmt.Errorf("no execution method found for worker %s", worker.GetName())
	}

	methodType := reflect.TypeOf(method)

	if methodType.NumIn() != 1 {
		return fmt.Errorf("execution method %s must have exactly one parameter, but found %d",
			methodType.Name(), methodType.NumIn())
	}

	uc.workerByName[worker.GetName()] = true
	uc.Workers = append(uc.Workers, *worker)

	workers := []workers.Worker{*worker}
	uc.registrationClient.AddWorkers(workers)
	return nil
}

func (uc *UnmeshedClient) RegisterWorker(worker *workers.Worker) error {
	return uc.registerWorker(worker)
}

func (uc *UnmeshedClient) RegisterWorkers(workers []*workers.Worker) error {
	for _, worker := range workers {
		if err := uc.registerWorker(worker); err != nil {
			return err
		}
	}
	return nil
}

func (uc *UnmeshedClient) Stop() {
	uc.stopPolling.Store(true)
	uc.stopOnce.Do(func() {
		close(uc.done)
	})
}

func (uc *UnmeshedClient) RunProcessSyncWithDefaultTimeout(processRequestData *common.ProcessRequestData) (*common.ProcessData, error) {
	return uc.processClient.RunProcessSync(processRequestData, 0)
}

func (uc *UnmeshedClient) RunProcessSync(processRequestData *common.ProcessRequestData, processTimeoutSeconds int) (*common.ProcessData, error) {
	return uc.processClient.RunProcessSync(processRequestData, processTimeoutSeconds)
}

func (uc *UnmeshedClient) RunProcessAsync(processRequestData *common.ProcessRequestData) (*common.ProcessData, error) {
	return uc.processClient.RunProcessAsync(processRequestData)
}

func (uc *UnmeshedClient) GetProcessData(processID int64, includeSteps bool, hideLargeValues bool) (*common.ProcessData, error) {
	return uc.processClient.GetProcessData(processID, includeSteps, hideLargeValues)
}

func (uc *UnmeshedClient) GetStepData(stepID int64) (*common.StepData, error) {
	return uc.processClient.GetStepData(stepID)
}

func (uc *UnmeshedClient) SearchProcessExecutions(params *common.ProcessSearchRequest) ([]*common.ProcessData, error) {
	return uc.processClient.SearchProcessExecutions(params)
}

func (uc *UnmeshedClient) InvokeAPIMappingGet(endpoint string, id string, correlationID string, apiCallType common.ApiCallType) (map[string]interface{}, error) {
	return uc.processClient.InvokeAPIMappingGet(endpoint, id, correlationID, apiCallType)
}

func (uc *UnmeshedClient) InvokeAPIMappingPost(endpoint string, input map[string]interface{}, id string, correlationID string, apiCallType common.ApiCallType) (map[string]interface{}, error) {
	return uc.processClient.InvokeAPIMappingPost(endpoint, input, id, correlationID, apiCallType)
}

func (uc *UnmeshedClient) BulkTerminate(processIDs []int64, reason string) (*common.ProcessActionResponseData, error) {
	return uc.processClient.BulkTerminate(processIDs, reason)
}

func (uc *UnmeshedClient) BulkResume(processIDs []int64) (*common.ProcessActionResponseData, error) {
	return uc.processClient.BulkResume(processIDs)
}

func (uc *UnmeshedClient) BulkReviewed(processIDs []int64, reason string) (*common.ProcessActionResponseData, error) {
	return uc.processClient.BulkReviewed(processIDs, reason)
}

func (uc *UnmeshedClient) Rerun(processID int64, version int) (*common.ProcessData, error) {
	return uc.processClient.Rerun(processID, version)
}

func (uc *UnmeshedClient) DoneChan() <-chan struct{} {
	return uc.done
}
