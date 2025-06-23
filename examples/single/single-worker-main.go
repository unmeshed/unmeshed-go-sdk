package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/main"
	apis2 "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

var unmeshedClient *apis.UnmeshedClient

// TestWorker is a simple worker function that returns a test message
func TestWorker(data map[string]interface{}) string {
	fmt.Println("TestWorker running with data:", data)
	time.Sleep(1 * time.Second)
	return "Hello from test-worker!"
}

func main() {
	worker := apis2.NewWorker(TestWorker, "test-worker")

	clientConfig := configs.NewClientConfig()
	clientConfig.SetClientID("<< Client Id >>")
	clientConfig.SetAuthToken("<< Auth Token >>")
	clientConfig.SetDelayMillis(40)
	clientConfig.SetPort(8080)
	clientConfig.SetWorkRequestBatchSize(500)
	clientConfig.SetBaseURL("http://localhost")
	clientConfig.SetStepTimeoutMillis(36000000)
	clientConfig.SetMaxWorkers(100)

	var err error
	unmeshedClient, err = apis.NewUnmeshedClient(clientConfig)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	unmeshedClient.RegisterWorker(worker)

	done := make(chan struct{})

	go func() {
		unmeshedClient.Start()
		close(done)
	}()

	fmt.Println("Single worker client started. Press Ctrl+C to stop...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("\nReceived shutdown signal. Stopping client...")
	case <-done:
		fmt.Println("Client finished execution")
	}
}
