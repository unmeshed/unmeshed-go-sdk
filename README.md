# Unmeshed Go SDK

Welcome to the **Unmeshed Go SDK**! 🚀  
This SDK lets you easily register and run workers that process jobs using the [Unmeshed](https://unmeshed.io) platform.

This README will walk you through everything you need: from installing the SDK, configuring your client, to writing and registering your own workers.

---

## Table of Contents

- [Installation](#installation)
- [Quickstart Example](#quickstart-example)
- [Configuration](#configuration)
- [Writing & Registering Workers](#writing--registering-workers)
- [Running the Client](#running-the-client)
- [Full Example](#full-example)

---

## Installation

First things first, let's get the SDK into your project:
```
go mod init

go get github.com/unmeshed/unmeshed-go-sdk
```
---

## Quickstart Example

Here's a super simple example to get you up and running:

```
package main

import (
    apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/main"
    apis2 "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
    "github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

func MyWorker(data map[string]interface{}) string {
    return "Hello, Unmeshed!"
}

func main() {
    worker := apis2.NewWorker(MyWorker, "my-worker")
    cfg := configs.NewClientConfig()
    cfg.SetClientID("your-client-id")
    cfg.SetAuthToken("your-auth-token")
    cfg.SetPort(8080)
    cfg.SetBaseURL("http://localhost")      // set your URL

    client := apis.NewUnmeshedClient(cfg)
    client.RegisterWorker(worker)
    client.Start()
}
```

---

## Configuration

Before you can start processing jobs, you'll need to configure your Unmeshed client. The config is flexible and lets you tweak things like:

- **Client ID & Auth Token:** For authentication (can be fetched from Unmeshed Instance).
- **Port:** (Optional for HTTPS) The port where your unmeshed server is running.
- **Base URL:** URL where your Unmeshed server is running.
- **Batch Size, Delays, Timeouts, Max Workers:** For tuning performance.

Example:
```
    cfg := configs.NewClientConfig()
    cfg.SetClientID("your-client-id")
    cfg.SetAuthToken("your-auth-token")
    cfg.SetPort(8080)
    cfg.SetBaseURL("http://localhost")      // set your URL
    cfg.SetWorkRequestBatchSize(50)
    cfg.SetInitialDelayMillis(50)
    cfg.SetStepTimeoutMillis(36000000)
    cfg.SetMaxWorkers(20)
```
---

## Writing & Registering Workers

A **worker** is just a Go function that takes some input and returns a result. You can use different input/output types, and even return multiple values!

Here are some examples:

```
// Sums all the values in the map
func Sum(data map[string]int) int {
    total := 0
    for _, v := range data {
    total += v
}
  return total
}
```

```
// Returns an error
func FailExample(data map[string]interface{}) error {
    return errors.New("failing worker")
}
```

To reschedule a worker with same iteration use return type of function as StepResult and
use KeepRunning and RescheduleAfterSeconds fields to control behaviour whether to complete or Reschedule
```
var GlobalCounter int = 0
reschedule it
func RescheduleWorkerExample(data map[string]interface{}) *common.StepResult {
	fmt.Println("Reschedule worker running with counter", GlobalCounter)
	if GlobalCounter > 5 {
		stepResult := common.NewStepResult("Response after 5 iterations")
		stepResult.KeepRunning = false
		stepResult.RescheduleAfterSeconds = 0
		GlobalCounter = 0
		return stepResult

	}
	GlobalCounter++
	stepResult := common.NewStepResult("Rescheduling the worker")
	stepResult.KeepRunning = true
	stepResult.RescheduleAfterSeconds = 2
	return stepResult
}
```

```
// Returns a list of strings
func ListExample(data map[string]interface{}) []string {
    return []string{"hello", "world"}
}
```

```
// Multiple outputs
func MultiOutputExample(data map[string]interface{}) (string, int) {
    return "Processed Successfully", len(data)
}
```

Register your workers like this:

```
import apis2 "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"

workers := []*apis2.Worker{
    apis2.NewWorker(RescheduleWorkerExample, "reschedule-worker)
    apis2.NewWorker(Sum, "sum"),
    apis2.NewWorker(FailExample, "fail-example"),
    apis2.NewWorker(ListExample, "list-example"),
    apis2.NewWorker(MultiOutputExample, "multi-output-example"),
}

client.RegisterWorkers(workers)

// You can also register workers one by one:
client.RegisterWorker(apis2.NewWorker(MyWorker, "my-worker"))
```

---

## Running the Client

Once everything's set up, just call:
client.Start()

The client will start polling for jobs and dispatching them to your workers. It's fully async and runs in the background.

---

## Full Example

Here's a more complete example that demonstrates all the available operations:

```go
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

// Example worker functions
func Sum(data map[string]int) int {
    sum := 0
    for _, v := range data {
        sum += v
    }
    return sum
}

func FailExample(data map[string]interface{}) error {
    return errors.New("failing worker")
}

func ListExample(data map[string]interface{}) []string {
    return []string{"hello", "world"}
}

func DelayedResponse(data map[string]interface{}) string {
    time.Sleep(10 * time.Second)
    return "Response after 10 seconds delay"
}

func main() {
    // Configure the client
    cfg := configs.NewClientConfig()
    cfg.SetClientID("your-client-id")
    cfg.SetAuthToken("your-auth-token")
    cfg.SetPort(8080)
    cfg.SetBaseURL("http://localhost")
    cfg.SetWorkRequestBatchSize(50)
    cfg.SetInitialDelayMillis(50)
    cfg.SetStepTimeoutMillis(36000000)
    cfg.SetMaxWorkers(20)

    // Create and configure the client
    client := apis.NewUnmeshedClient(cfg)

    // Register workers
    workers := []*apis2.Worker{
        apis2.NewWorker(Sum, "sum"),
        apis2.NewWorker(FailExample, "fail-example"),
        apis2.NewWorker(ListExample, "list-example"),
        apis2.NewWorker(DelayedResponse, "delayed-response"),
    }
    client.RegisterWorkers(workers)

    // Start the client in a goroutine
    done := make(chan struct{})
    go func() {
        client.Start()
        close(done)
    }()

    // Example 1: Create and run a process synchronously
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

    processData1, err := client.RunProcessSyncWithDefaultTimeout(processRequest)
    if err != nil {
        fmt.Printf("Error running process sync: %v\n", err)
        return
    }
    fmt.Printf("Sync execution returned: %+v\n", processData1)

    // Example 2: Run process asynchronously
    processData2, err := client.RunProcessAsync(processRequest)
    if err != nil {
        fmt.Printf("Error running process async: %v\n", err)
        return
    }
    fmt.Printf("Async execution returned: %+v\n", processData2)

    // Example 3: Get process data without steps
    processData1Retrieved1, err := client.GetProcessData(processData1.ProcessID, false)
    if err != nil {
        fmt.Printf("Error getting process data: %v\n", err)
        return
    }
    fmt.Printf("Process data without steps: %+v\n", processData1Retrieved1)

    // Example 4: Get process data with steps
    processData1Retrieved2, err := client.GetProcessData(processData1.ProcessID, true)
    if err != nil {
        fmt.Printf("Error getting process data with steps: %v\n", err)
        return
    }
    fmt.Printf("Process data with steps: %+v\n", processData1Retrieved2)

    // Example 5: Get step data
    if len(processData1Retrieved2.StepRecords) > 0 {
        stepData1, err := client.GetStepData(processData1Retrieved2.StepRecords[0].StepID)
        if err != nil {
            fmt.Printf("Error getting step data: %v\n", err)
            return
        }
        fmt.Printf("Step data: %+v\n", stepData1)
    }

    // Example 6: Search process executions
    searchNamespace := "default"
    searchNames := []string{"test_process"}
    processSearchRequest := &common.ProcessSearchRequest{
        Names:     searchNames,
        Limit:     20,
        Namespace: &searchNamespace,
    }
    processesSearchResults, err := client.SearchProcessExecutions(processSearchRequest)
    if err != nil {
        fmt.Printf("Error searching process executions: %v\n", err)
        return
    }
    fmt.Printf("Search returned %d results\n", len(processesSearchResults))

    // Example 7: Rerun process
    rerunProcessData, err := client.Rerun(processData1.ProcessID, 1)
    if err != nil {
        fmt.Printf("Error rerunning process: %v\n", err)
        return
    }
    fmt.Printf("Rerun process data: %+v\n", rerunProcessData)

    // Example 8: Bulk terminate processes
    actionResponse, err := client.BulkTerminate([]int64{processData1.ProcessID, 1, 2}, "Terminating processes")
    if err != nil {
        fmt.Printf("Error bulk terminating processes: %v\n", err)
        return
    }
    fmt.Printf("Bulk terminate response: %+v\n", actionResponse.Details)

    // Example 9: Bulk resume processes
    actionResponse, err = client.BulkResume([]int64{processData1.ProcessID, 1, 2})
    if err != nil {
        fmt.Printf("Error bulk resuming processes: %v\n", err)
        return
    }
    fmt.Printf("Bulk resume response: %+v\n", actionResponse.Details)

    // Example 10: Bulk review processes
    actionResponse, err = client.BulkReviewed([]int64{processData1.ProcessID, 1, 2}, "Reviewing processes")
    if err != nil {
        fmt.Printf("Error bulk reviewing processes: %v\n", err)
        return
    }
    fmt.Printf("Bulk review response: %+v\n", actionResponse.Details)

    // Example 11: Invoke API mapping GET
    response, err := client.InvokeAPIMappingGet(
        "test_process_endpoint",
        "req_id--1",
        "correl_id--1",
        common.ApiCallTypeSync,
    )
    if err != nil {
        fmt.Printf("Error invoking API mapping GET: %v\n", err)
        return
    }
    fmt.Printf("API mapping GET response: %+v\n", response)

    // Example 12: Invoke API mapping POST
    response, err = client.InvokeAPIMappingPost(
        "test_process_endpoint",
        map[string]interface{}{"test": "value"},
        "req_id--1",
        "correl_id--1",
        common.ApiCallTypeSync,
    )
    if err != nil {
        fmt.Printf("Error invoking API mapping POST: %v\n", err)
        return
    }
    fmt.Printf("API mapping POST response: %+v\n", response)

    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    select {
    case <-sigChan:
        fmt.Println("\nReceived shutdown signal. Stopping client...")
    case <-done:
        fmt.Println("Client finished execution")
    }
}
```

This example demonstrates:
1. Client configuration and worker registration
2. Running processes synchronously and asynchronously
3. Retrieving process data with and without steps
4. Getting step data
5. Searching process executions
6. Rerunning processes
7. Bulk operations (terminate, resume, review)
8. API mapping invocations (GET and POST)
9. Graceful shutdown handling

Each operation includes proper error handling and logging of results. The client is started in a goroutine to allow for concurrent operations, and the program handles graceful shutdown through signal handling.

---
