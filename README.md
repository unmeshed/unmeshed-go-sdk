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

Here’s a more complete example :

```
package main

import (
    "errors"
    "fmt"
    apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/main"
    apis2 "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
    "github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

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

func main() {
    workers := []*apis2.Worker{
        apis2.NewWorker(Sum, "sum"),
        apis2.NewWorker(FailExample, "fail-example"),
        apis2.NewWorker(ListExample, "list-example"),
    }

    cfg := configs.NewClientConfig()
    cfg.SetClientID("your-client-id")
    cfg.SetAuthToken("your-auth-token")
    cfg.SetPort(8080)
    cfg.SetBaseURL("http://localhost")

    client := apis.NewUnmeshedClient(cfg)
    client.RegisterWorkers(workers)
    client.Start()
}
```

---
