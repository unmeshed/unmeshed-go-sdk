package main

import (
	"errors"
	"fmt"
	"os"
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
	clientConfig.SetWorkRequestBatchSize(50)
	clientConfig.SetBaseURL("http://localhost")
	clientConfig.SetInitialDelayMillis(50)
	clientConfig.SetStepTimeoutMillis(36000000)
	clientConfig.SetMaxWorkers(20)
	unmeshedClient := apis.NewUnmeshedClient(clientConfig)
	unmeshedClient.RegisterWorkers(workerList)

	worker := apis2.NewWorker(ManuallyRegisteredWorker, "manually-registered-worker")
	unmeshedClient.RegisterWorker(worker)

	// Start the client - it will run in background
	unmeshedClient.Start()
}
