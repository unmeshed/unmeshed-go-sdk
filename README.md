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
- [Process Definition Management](#process-definition-management)
- [Full Example](#full-example)

---

## Installation

First things first, let's get the SDK into your project:

```bash
go mod init

go get github.com/unmeshed/unmeshed-go-sdk
```

---

## Quickstart Example

Here's a super simple example to get you up and running:

```go
package main

import (
    "fmt"
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

    unmeshedClient, err := apis.NewUnmeshedClient(cfg)
    if err != nil {
        fmt.Printf("Error creating client: %v\n", err)
        return
    }
    unmeshedClient.RegisterWorker(worker)
    unmeshedClient.Start()
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

```go
cfg := configs.NewClientConfig()
cfg.SetClientID("your-client-id")
cfg.SetAuthToken("your-auth-token")
cfg.SetPort(8080)
cfg.SetBaseURL("http://localhost")      // set your URL
cfg.SetWorkRequestBatchSize(50)
cfg.SetDelayMillis(100) //Delay between poll
cfg.SetSubmitClientSleepIntervalMillis(100) // Delay between submit attempts
cfg.SetStepTimeoutMillis(36000000)
cfg.SetMaxWorkers(20)
```

---

## Writing & Registering Workers

A **worker** is just a Go function that takes some input and returns a result. You can use different input/output types, and even return multiple values!

Here are some examples:

```go
// Sums all the values in the map
func Sum(data map[string]int) int {
    total := 0
    for _, v := range data {
        total += v
    }
    return total
}
```

```go
// Returns an error
func FailExample(data map[string]interface{}) error {
    return errors.New("failing worker")
}
```

Getting currently executed workRequest by client:

```go
workRequest := unmeshedClient.GetCurrentWorkRequest()
fmt.Println(workRequest)
```

Display/Hide Large values as part of output during process search.

Use flag `hideLargeValues` as part of `GetProcessData(processId, includeSteps, hideLargeValues)` to include/exclude large output payloads in execution response:

```go
processData1Retrieved1, err := unmeshedClient.GetProcessData(processData1.ProcessID, false, false)
if err != nil {
    fmt.Printf("Error getting process data: %v\n", err)
    return
}
fmt.Printf("Retrieving process %d returned %+v\n", processData1.ProcessID, processData1Retrieved1)
fmt.Printf("Since the flag to include steps was false the steps was not returned: %d\n", len(processData1Retrieved1.StepRecords))
```

To reschedule a worker with same iteration use return type of function as `StepResult` and use `KeepRunning` and `RescheduleAfterSeconds` fields to control behaviour whether to complete or reschedule it:

```go
var GlobalCounter int = 0

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

```go
// Returns a list of strings
func ListExample(data map[string]interface{}) []string {
    return []string{"hello", "world"}
}
```

```go
// Multiple outputs
func MultiOutputExample(data map[string]interface{}) (string, int) {
    return "Processed Successfully", len(data)
}
```

Register your workers like this:

```go
import apis2 "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"

workers := []*apis2.Worker{
    apis2.NewWorker(RescheduleWorkerExample, "reschedule-worker"),
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

```go
client.Start()
```

The client will start polling for jobs and dispatching them to your workers. It's fully async and runs in the background.

---

## 🧩 Process Definition Management

You can manage process definitions directly from the SDK — create, update, fetch, and delete definitions programmatically.

### 🏗️ Create a Process Definition

```go
noop1 := &common.StepDefinition{
    Name:        StringPtr("noop1"),
    Ref:         StringPtr("noop1"),
    Description: StringPtr("Test noop 1"),
    Type:        StepTypePtr(common.NOOP),
    Input:       map[string]interface{}{"key1": "val1"},
}

noop2 := &common.StepDefinition{
    Name:        StringPtr("noop2"),
    Ref:         StringPtr("noop2"),
    Description: StringPtr("Test noop 2"),
    Type:        StepTypePtr(common.NOOP),
    Input:       map[string]interface{}{"key1": "val1"},
}

processDef := &common.ProcessDefinition{
    Name:        StringPtr("test-process"),
    Version:     IntPtr(1),
    Namespace:   StringPtr("default"),
    Description: StringPtr("Testing Process"),
    Type:        StringPtr("API_ORCHESTRATION"),
    Steps:       []*common.StepDefinition{noop1, noop2},
}

created, err := client.CreateNewProcessDefinition(processDef)
if err != nil {
    fmt.Printf("Error creating process definition: %v\n", err)
    return
}
fmt.Printf("Created process definition: %+v\n", created)
```

### ✏️ Update a Process Definition

```go
noop3 := &common.StepDefinition{
    Name:        StringPtr("noop3"),
    Ref:         StringPtr("noop3"),
    Description: StringPtr("Test noop 3"),
    Type:        StepTypePtr(common.NOOP),
    Input:       map[string]interface{}{"key1": "val1"},
}

updatedDef := &common.ProcessDefinition{
    Name:        StringPtr("test-process"),
    Version:     IntPtr(2),
    Namespace:   StringPtr("default"),
    Description: StringPtr("Testing Process Updated"),
    Type:        StringPtr("API_ORCHESTRATION"),
    Steps:       []*common.StepDefinition{noop1, noop2, noop3},
}

updated, err := client.UpdateProcessDefinition(updatedDef)
if err != nil {
    fmt.Printf("Error updating process definition: %v\n", err)
    return
}
fmt.Printf("Updated process definition: %+v\n", updated)
```

### 🗑️ Delete Process Definitions

```go
defs, err := client.GetAllProcessDefinitions()
if err != nil {
    fmt.Printf("Error fetching definitions: %v\n", err)
    return
}

var toDelete []*common.ProcessDefinition
for _, pd := range defs {
    if *pd.Name == "test-process" && *pd.Namespace == "default" {
        toDelete = append(toDelete, pd)
    }
}

if len(toDelete) > 0 {
    response, err := client.DeleteProcessDefinitions(toDelete, false) // false = delete all versions
    if err != nil {
        fmt.Printf("Error deleting definitions: %v\n", err)
        return
    }
    fmt.Printf("Delete response: %+v\n", response)
}
```

Delete only a specific version:

```go
version := 1
var version1Defs []*common.ProcessDefinition
for _, pd := range defs {
    if *pd.Name == "test-process" && *pd.Namespace == "default" && *pd.Version == version {
        version1Defs = append(version1Defs, pd)
    }
}

if len(version1Defs) > 0 {
    response, err := client.DeleteProcessDefinitions(version1Defs, true) // true = delete only version 1
    if err != nil {
        fmt.Printf("Error deleting version: %v\n", err)
        return
    }
    fmt.Printf("Delete version response: %+v\n", response)
}
```

### 🔍 Get Process Definitions

```go
// Get latest version
latest, err := client.GetProcessDefinitionLatestOrVersion("default", "test-process", nil)
if err != nil {
    fmt.Printf("Error getting latest definition: %v\n", err)
    return
}
fmt.Printf("Latest definition: name=%s, version=%d\n", *latest.Name, *latest.Version)

// Get specific version
version := 1
specific, err := client.GetProcessDefinitionLatestOrVersion("default", "test-process", &version)
if err != nil {
    fmt.Printf("Error getting specific definition: %v\n", err)
    return
}
fmt.Printf("Specific version: name=%s, version=%d\n", *specific.Name, *specific.Version)
```

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

// Helper functions for pointer conversion
func StringPtr(s string) *string {
    return &s
}

func IntPtr(i int) *int {
    return &i
}

func StepTypePtr(st common.StepType) *common.StepType {
    return &st
}

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
    cfg.SetDelayMillis(100) //Delay between poll
    cfg.SetStepTimeoutMillis(36000000)
    cfg.SetSubmitClientSleepIntervalMillis(100)
    cfg.SetMaxWorkers(20)

    // Create and configure the client
    unmeshedClient, err := apis.NewUnmeshedClient(cfg)
    if err != nil {
        fmt.Printf("Error creating client: %v\n", err)
        return
    }

    // Register workers
    workers := []*apis2.Worker{
        apis2.NewWorker(Sum, "sum"),
        apis2.NewWorker(FailExample, "fail-example"),
        apis2.NewWorker(ListExample, "list-example"),
        apis2.NewWorker(DelayedResponse, "delayed-response"),
    }
    unmeshedClient.RegisterWorkers(workers)

    // Start the client in a goroutine
    done := make(chan struct{})
    go func() {
        unmeshedClient.Start()
        close(done)
    }()

    // ========================================
    // PROCESS DEFINITION MANAGEMENT EXAMPLES
    // ========================================

    // Example 1: Create a new process definition
    fmt.Println("\n=== Creating Process Definition ===")
    noop1 := &common.StepDefinition{
        Name:        StringPtr("noop1"),
        Ref:         StringPtr("noop1"),
        Description: StringPtr("Test noop 1"),
        Type:        StepTypePtr(common.NOOP),
        Input:       map[string]interface{}{"key1": "val1"},
    }

    noop2 := &common.StepDefinition{
        Name:        StringPtr("noop2"),
        Ref:         StringPtr("noop2"),
        Description: StringPtr("Test noop 2"),
        Type:        StepTypePtr(common.NOOP),
        Input:       map[string]interface{}{"key1": "val1"},
    }

    processDef := &common.ProcessDefinition{
        Name:        StringPtr("test-process"),
        Version:     IntPtr(1),
        Namespace:   StringPtr("default"),
        Description: StringPtr("Testing Process"),
        Type:        StringPtr("API_ORCHESTRATION"),
        Steps:       []*common.StepDefinition{noop1, noop2},
    }

    created, err := unmeshedClient.CreateNewProcessDefinition(processDef)
    if err != nil {
        fmt.Printf("Error creating process definition: %v\n", err)
    } else {
        fmt.Printf("Created process definition: %+v\n", created)
    }

    // Example 2: Get latest process definition
    fmt.Println("\n=== Getting Latest Process Definition ===")
    latest, err := unmeshedClient.GetProcessDefinitionLatestOrVersion("default", "test-process", nil)
    if err != nil {
        fmt.Printf("Error getting latest definition: %v\n", err)
    } else {
        fmt.Printf("Latest definition: name=%s, version=%d\n", *latest.Name, *latest.Version)
    }

    // Example 3: Update process definition
    fmt.Println("\n=== Updating Process Definition ===")
    noop3 := &common.StepDefinition{
        Name:        StringPtr("noop3"),
        Ref:         StringPtr("noop3"),
        Description: StringPtr("Test noop 3"),
        Type:        StepTypePtr(common.NOOP),
        Input:       map[string]interface{}{"key1": "val1"},
    }

    updatedDef := &common.ProcessDefinition{
        Name:        StringPtr("test-process"),
        Version:     IntPtr(2),
        Namespace:   StringPtr("default"),
        Description: StringPtr("Testing Process Updated"),
        Type:        StringPtr("API_ORCHESTRATION"),
        Steps:       []*common.StepDefinition{noop1, noop2, noop3},
    }

    updated, err := unmeshedClient.UpdateProcessDefinition(updatedDef)
    if err != nil {
        fmt.Printf("Error updating process definition: %v\n", err)
    } else {
        fmt.Printf("Updated process definition: %+v\n", updated)
    }

    // Example 4: Get specific version of process definition
    fmt.Println("\n=== Getting Specific Version of Process Definition ===")
    version := 1
    specific, err := unmeshedClient.GetProcessDefinitionLatestOrVersion("default", "test-process", &version)
    if err != nil {
        fmt.Printf("Error getting specific definition: %v\n", err)
    } else {
        fmt.Printf("Specific version: name=%s, version=%d\n", *specific.Name, *specific.Version)
    }

    // Example 5: Get all process definitions
    fmt.Println("\n=== Getting All Process Definitions ===")
    allDefs, err := unmeshedClient.GetAllProcessDefinitions()
    if err != nil {
        fmt.Printf("Error fetching all definitions: %v\n", err)
    } else {
        fmt.Printf("Total process definitions: %d\n", len(allDefs))
    }

    // ========================================
    // PROCESS EXECUTION EXAMPLES
    // ========================================

    // Example 6: Create and run a process synchronously
    fmt.Println("\n=== Running Process Synchronously ===")
    namespace := "default"
    name := "test_process"
    processVersion := 1
    requestID := "req001"
    correlationID := "corr001"
    processRequest := &common.ProcessRequestData{
        Namespace:     &namespace,
        Name:          &name,
        Version:       &processVersion,
        RequestID:     &requestID,
        CorrelationID: &correlationID,
        Input: map[string]interface{}{
            "test1": "value",
            "test2": 100,
            "test3": 100.0,
        },
    }

    processData1, err := unmeshedClient.RunProcessSyncWithDefaultTimeout(processRequest)
    if err != nil {
        fmt.Printf("Error running process sync: %v\n", err)
    } else {
        fmt.Printf("Sync execution returned: %+v\n", processData1)
    }

    // Example 7: Run process asynchronously
    fmt.Println("\n=== Running Process Asynchronously ===")
    processData2, err := unmeshedClient.RunProcessAsync(processRequest)
    if err != nil {
        fmt.Printf("Error running process async: %v\n", err)
    } else {
        fmt.Printf("Async execution returned: %+v\n", processData2)
    }

    // Example 8: Get process data without steps
    fmt.Println("\n=== Getting Process Data Without Steps ===")
    if processData1 != nil {
        processData1Retrieved1, err := unmeshedClient.GetProcessData(processData1.ProcessID, false, false)
        if err != nil {
            fmt.Printf("Error getting process data: %v\n", err)
        } else {
            fmt.Printf("Process data without steps: %+v\n", processData1Retrieved1)
        }

        // Example 9: Get process data with steps
        fmt.Println("\n=== Getting Process Data With Steps ===")
        processData1Retrieved2, err := unmeshedClient.GetProcessData(processData1.ProcessID, true, false)
        if err != nil {
            fmt.Printf("Error getting process data with steps: %v\n", err)
        } else {
            fmt.Printf("Process data with steps: %+v\n", processData1Retrieved2)

            // Example 10: Get step data
            fmt.Println("\n=== Getting Step Data ===")
            if len(processData1Retrieved2.StepRecords) > 0 {
                stepData1, err := unmeshedClient.GetStepData(processData1Retrieved2.StepRecords[0].StepID)
                if err != nil {
                    fmt.Printf("Error getting step data: %v\n", err)
                } else {
                    fmt.Printf("Step data: %+v\n", stepData1)
                }
            }
        }
    }

    // Example 11: Search process executions
    fmt.Println("\n=== Searching Process Executions ===")
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
    } else {
        fmt.Printf("Search returned %d results\n", len(processesSearchResults))
    }

    // Example 12: Rerun process
    fmt.Println("\n=== Rerunning Process ===")
    if processData1 != nil {
        rerunProcessData, err := unmeshedClient.Rerun(processData1.ProcessID, 1)
        if err != nil {
            fmt.Printf("Error rerunning process: %v\n", err)
        } else {
            fmt.Printf("Rerun process data: %+v\n", rerunProcessData)
        }
    }

    // Example 13: Bulk terminate processes
    fmt.Println("\n=== Bulk Terminating Processes ===")
    if processData1 != nil {
        actionResponse, err := unmeshedClient.BulkTerminate([]int64{processData1.ProcessID, 1, 2}, "Terminating processes")
        if err != nil {
            fmt.Printf("Error bulk terminating processes: %v\n", err)
        } else {
            fmt.Printf("Bulk terminate response: %+v\n", actionResponse.Details)
        }
    }

    // Example 14: Bulk resume processes
    fmt.Println("\n=== Bulk Resuming Processes ===")
    if processData1 != nil {
        actionResponse, err := unmeshedClient.BulkResume([]int64{processData1.ProcessID, 1, 2})
        if err != nil {
            fmt.Printf("Error bulk resuming processes: %v\n", err)
        } else {
            fmt.Printf("Bulk resume response: %+v\n", actionResponse.Details)
        }
    }

    // Example 15: Bulk review processes
    fmt.Println("\n=== Bulk Reviewing Processes ===")
    if processData1 != nil {
        actionResponse, err := unmeshedClient.BulkReviewed([]int64{processData1.ProcessID, 1, 2}, "Reviewing processes")
        if err != nil {
            fmt.Printf("Error bulk reviewing processes: %v\n", err)
        } else {
            fmt.Printf("Bulk review response: %+v\n", actionResponse.Details)
        }
    }

    // Example 16: Invoke API mapping GET
    fmt.Println("\n=== Invoking API Mapping GET ===")
    response, err := unmeshedClient.InvokeAPIMappingGet(
        "test_process_endpoint",
        "req_id--1",
        "correl_id--1",
        common.ApiCallTypeSync,
    )
    if err != nil {
        fmt.Printf("Error invoking API mapping GET: %v\n", err)
    } else {
        fmt.Printf("API mapping GET response: %+v\n", response)
    }

    // Example 17: Invoke API mapping POST
    fmt.Println("\n=== Invoking API Mapping POST ===")
    response, err = unmeshedClient.InvokeAPIMappingPost(
        "test_process_endpoint",
        map[string]interface{}{"test": "value"},
        "req_id--1",
        "correl_id--1",
        common.ApiCallTypeSync,
    )
    if err != nil {
        fmt.Printf("Error invoking API mapping POST: %v\n", err)
    } else {
        fmt.Printf("API mapping POST response: %+v\n", response)
    }

    // ========================================
    // PROCESS DEFINITION DELETION EXAMPLES
    // ========================================

    // Example 18: Delete specific version of process definition
    fmt.Println("\n=== Deleting Specific Version of Process Definition ===")
    defs, err := unmeshedClient.GetAllProcessDefinitions()
    if err != nil {
        fmt.Printf("Error fetching definitions: %v\n", err)
    } else {
        deleteVersion := 1
        var version1Defs []*common.ProcessDefinition
        for _, pd := range defs {
            if *pd.Name == "test-process" && *pd.Namespace == "default" && *pd.Version == deleteVersion {
                version1Defs = append(version1Defs, pd)
            }
        }

        if len(version1Defs) > 0 {
            deleteResponse, err := unmeshedClient.DeleteProcessDefinitions(version1Defs, true) // true = delete only version 1
            if err != nil {
                fmt.Printf("Error deleting version: %v\n", err)
            } else {
                fmt.Printf("Delete version response: %+v\n", deleteResponse)
            }
        }
    }

    // Example 19: Delete all versions of process definition
    fmt.Println("\n=== Deleting All Versions of Process Definition ===")
    defs, err = unmeshedClient.GetAllProcessDefinitions()
    if err != nil {
        fmt.Printf("Error fetching definitions: %v\n", err)
    } else {
        var toDelete []*common.ProcessDefinition
        for _, pd := range defs {
            if *pd.Name == "test-process" && *pd.Namespace == "default" {
                toDelete = append(toDelete, pd)
            }
        }

        if len(toDelete) > 0 {
            deleteResponse, err := unmeshedClient.DeleteProcessDefinitions(toDelete, false) // false = delete all versions
            if err != nil {
                fmt.Printf("Error deleting definitions: %v\n", err)
            } else {
                fmt.Printf("Delete response: %+v\n", deleteResponse)
            }
        }
    }

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

This comprehensive example demonstrates:

1. **Process Definition Management:**
   - Creating new process definitions
   - Getting latest version of a definition
   - Updating process definitions
   - Getting specific versions
   - Getting all process definitions
   - Deleting specific versions
   - Deleting all versions

2. **Process Execution:**
   - Running processes synchronously and asynchronously
   - Retrieving process data with and without steps
   - Getting step data
   - Searching process executions
   - Rerunning processes

3. **Bulk Operations:**
   - Bulk terminate
   - Bulk resume
   - Bulk review

4. **API Mapping:**
   - GET and POST invocations

5. **Error Handling & Graceful Shutdown:**
   - Proper error handling for all operations
   - Signal handling for graceful shutdown

Each operation includes descriptive headers and error handling to make it easy to understand and debug.

---