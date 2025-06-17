package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/main"
	apis2 "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

var GlobalCounter int = 0

type MathOperations struct{}

func (m *MathOperations) Sum(data map[string]int) int {
	sum := 0
	for _, v := range data {
		sum += v
	}
	return sum
}

func RescheduleExample(data map[string]interface{}) *common.StepResult {
	fmt.Println("Reschedule worker running with counter", GlobalCounter)
	if GlobalCounter > 5 {
		stepResult := common.NewStepResult("Response after 5 iterations")
		stepResult.KeepRunning = false
		stepResult.RescheduleAfterSeconds = 0
		GlobalCounter = 0
		return stepResult

	}
	GlobalCounter++
	stepResult := common.NewStepResult("testing reschedule")
	stepResult.KeepRunning = true
	stepResult.RescheduleAfterSeconds = 2
	return stepResult
}

func FailExample(data map[string]interface{}) error {
	return errors.New("failing worker")
}

func ListExample(data map[string]interface{}) []string {
	result := make([]string, 0)
	result = append(result, "hello")
	result = append(result, "world")
	return result
}

func CountTrueValues(data map[int]bool) int {
	count := 0
	for _, v := range data {
		if v {
			count++
		}
	}
	return count
}

func MaxValue(data map[string]int) int {
	max := 0
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

func ProcessMap(data map[string]interface{}) map[string]interface{} {
	fmt.Println("ProcessMap running")
	return data
}

func LongestString(data map[string]string) string {
	longest := ""
	for _, v := range data {
		if len(v) > len(longest) {
			longest = v
		}
	}
	return longest
}

func KeyWithMaxValue(data map[string]int) string {
	maxKey := ""
	maxValue := 0
	for k, v := range data {
		if v > maxValue {
			maxKey = k
			maxValue = v
		}
	}
	return maxKey
}

type Person struct {
	Name   string                 `json:"name"`
	Age    int                    `json:"age"`
	Active bool                   `json:"active"`
	Score  float64                `json:"score"`
	Meta   map[string]interface{} `json:"meta"`
}

func GetPerson(p Person) Person {
	return p
}

func ReverseMap(data map[string]string) map[string]string {
	reversed := make(map[string]string)
	for k, v := range data {
		reversed[v] = k
	}
	return reversed
}

func FlattenValues(data map[string][]int) []int {
	var result []int
	for _, values := range data {
		result = append(result, values...)
	}
	return result
}

func ManuallyRegisteredWorker(data map[string]interface{}) string {
	return "Test"
}

func MultiOutputExample(data map[string]interface{}) (string, int) {
	message := "Processed Successfully"
	count := len(data)
	return message, count
}

func DelayedResponse(data map[string]interface{}) string {
	time.Sleep(10 * time.Second)
	return "Response after 10 seconds delay"
}

func main() {
	os.Setenv("DISABLE_RUNNING_WORKER_LOGS", "true")
	workerList := []*apis2.Worker{
		apis2.NewWorker(RescheduleExample, "reschedule-example"),
		apis2.NewWorker(DelayedResponse, "delayed-response"),
		apis2.NewWorker(MultiOutputExample, "multi-output-example"),
		apis2.NewWorker(ListExample, "list-example"),
		apis2.NewWorker(FailExample, "fail-example"),
		apis2.NewWorker(ProcessMap, "process-map"),
		apis2.NewWorker(GetPerson, "get_person"),
		apis2.NewWorker(CountTrueValues, "count_true_values"),
		apis2.NewWorker(MaxValue, "max_value"),
		apis2.NewWorker(LongestString, "longest_string"),
		apis2.NewWorker(KeyWithMaxValue, "key_with_max_value"),
		apis2.NewWorker(ReverseMap, "reverse_map"),
		apis2.NewWorker(FlattenValues, "flatten_values"),
	}

	clientConfig := configs.NewClientConfig()
	clientConfig.SetClientID("<< Client Id >>")
	clientConfig.SetAuthToken("<< Auth Token >>")
	clientConfig.SetPort(8080)
	clientConfig.SetWorkRequestBatchSize(200)
	clientConfig.SetBaseURL("http://localhost")
	clientConfig.SetStepTimeoutMillis(36000000)
	clientConfig.SetMaxWorkers(20)
	unmeshedClient := apis.NewUnmeshedClient(clientConfig)
	unmeshedClient.RegisterWorkers(workerList)

	worker := apis2.NewWorker(ManuallyRegisteredWorker, "manually-registered-worker")
	unmeshedClient.RegisterWorker(worker)

	done := make(chan struct{})

	///Start the client in goroutine
	go func() {
		unmeshedClient.Start()
		close(done)
	}()

	namespace := "default"
	name := "test_process"
	version := 1
	requestID := "req001"
	correlationID := "corr001"
	processRequest := &common.ProcessRequestData{
		Namespace:     &namespace,
		Name:          &name,
		Version:       &version,
		RequestID:     &requestID,
		CorrelationID: &correlationID,
		Input: map[string]interface{}{
			"test1": "value",
			"test2": 100,
			"test3": 100.0,
		},
	}

	// Run process synchronously
	processData1, err := unmeshedClient.RunProcessSyncWithDefaultTimeout(processRequest)
	if err != nil {
		fmt.Printf("Error running process sync: %v\n", err)
		return
	}
	fmt.Printf("Sync execution of process request %+v returned %+v\n", processRequest, processData1)

	// Run process asynchronously
	processData2, err := unmeshedClient.RunProcessAsync(processRequest)
	if err != nil {
		fmt.Printf("Error running process async: %v\n", err)
		return
	}
	fmt.Printf("Async execution of process request %+v returned %+v\n", processRequest, processData2)

	// Get process data without steps
	processData1Retrieved1, err := unmeshedClient.GetProcessData(processData1.ProcessID, false, false)
	if err != nil {
		fmt.Printf("Error getting process data: %v\n", err)
		return
	}
	fmt.Printf("Retrieving process %d returned %+v\n", processData1.ProcessID, processData1Retrieved1)
	fmt.Printf("Since the flag to include steps was false the steps was not returned: %d\n", len(processData1Retrieved1.StepRecords))

	// Get process data with steps
	processData1Retrieved2, err := unmeshedClient.GetProcessData(processData1.ProcessID, true, false)
	if err != nil {
		fmt.Printf("Error getting process data with steps: %v\n", err)
		return
	}
	fmt.Printf("Retrieving process %d returned %+v\n", processData1.ProcessID, processData1Retrieved2)
	fmt.Printf("Since the flag to include steps was true the steps was returned: %d\n", len(processData1Retrieved2.StepRecords))

	// Get step data
	if len(processData1Retrieved2.StepRecords) > 0 {
		stepData1, err := unmeshedClient.GetStepData(3950120)
		if err != nil {
			fmt.Printf("Error getting step data: %v\n", err)
			return
		}
		fmt.Printf("Retrieving step data %d returned %+v\n", stepData1.StepID, stepData1)
	}

	// Search process executions
	searchNamespace := "default"
	searchNames := []string{"test_process"}
	processSearchRequest := &common.ProcessSearchRequest{
		Names:     searchNames,
		Limit:     20,
		Namespace: &searchNamespace,
	}
	processesSearchResults, err := unmeshedClient.SearchProcessExecutions(processSearchRequest)
	if err != nil {
		fmt.Printf("Error searching process executions: %v\n", err)
		return
	}
	fmt.Printf("Search returned %d results\n", len(processesSearchResults))

	// Rerun process
	rerunProcessData, err := unmeshedClient.Rerun(processData1.ProcessID, 1)
	if err != nil {
		fmt.Printf("Error rerunning process: %v\n", err)
		return
	}
	fmt.Printf("Rerun of process %d returned %+v\n", processData1.ProcessID, rerunProcessData)

	// Bulk terminate processes
	actionResponse, err := unmeshedClient.BulkTerminate([]int64{3950142}, "Terminating processes")
	if err != nil {
		fmt.Printf("Error bulk terminating processes: %v\n", err)
		return
	}
	fmt.Printf("Bulk terminate of 1 process returned %+v\n", actionResponse.Details)

	// Bulk resume processes
	actionResponse, err = unmeshedClient.BulkResume([]int64{3950154})
	if err != nil {
		fmt.Printf("Error bulk resuming processes: %v\n", err)
		return
	}
	fmt.Printf("Bulk resume of 1 process returned %+v\n", actionResponse.Details)

	// Bulk review processes
	actionResponse, err = unmeshedClient.BulkReviewed([]int64{3950184}, "Reviewing processes")
	if err != nil {
		fmt.Printf("Error bulk reviewing processes: %v\n", err)
		return
	}
	fmt.Printf("Bulk review of 1 process returned %+v\n", actionResponse.Details)

	// Invoke API mapping GET
	response, err := unmeshedClient.InvokeAPIMappingGet(
		"test_process_endpoint/O0ZjQNLhc5JGjcLKVajV/w4VliQzkGnOafPk8V6AJ",
		"req_id--1",
		"correl_id--1",
		common.ApiCallTypeSync,
	)
	if err != nil {
		fmt.Printf("Error invoking API mapping GET: %v\n", err)
		return
	}
	fmt.Printf("API mapped endpoint invocation using GET returned %+v\n", response)

	// Invoke API mapping POST
	response, err = unmeshedClient.InvokeAPIMappingPost(
		"test_process_endpoint/O0ZjQNLhc5JGjcLKVajV/w4VliQzkGnOafPk8V6AJ",
		map[string]interface{}{"test": "value"},
		"req_id--1",
		"correl_id--1",
		common.ApiCallTypeSync,
	)
	if err != nil {
		fmt.Printf("Error invoking API mapping POST: %v\n", err)
		return
	}
	fmt.Printf("API mapped endpoint invocation using POST returned %+v\n", response)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("\nReceived shutdown signal. Stopping client...")
	case <-done:
		fmt.Println("Client finished execution")
	}
}
