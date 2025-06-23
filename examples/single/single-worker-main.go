package main

import (
	"fmt"
	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/main"
	apis2 "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var unmeshedClient *apis.UnmeshedClient

// TestWorker is a simple worker function that returns a test message
func TestWorker(data map[string]interface{}) string {
	fmt.Println("TestWorker running with data:", data)
	time.Sleep(1 * time.Second)
	return "Hello from test-worker!"
}

func main() {
	// Create a single worker
	worker := apis2.NewWorker(TestWorker, "test-worker")

	// Configure the client
	clientConfig := configs.NewClientConfig()
	clientConfig.SetClientID("0d9ea975-da46-4bff-b4b9-96be0f8a11b0")
	clientConfig.SetAuthToken("MAoL7MFfKagCFtEW2DUS")
	clientConfig.SetDelayMillis(40)
	clientConfig.SetPort(8080)
	clientConfig.SetWorkRequestBatchSize(500)
	clientConfig.SetBaseURL("http://localhost")
	clientConfig.SetStepTimeoutMillis(36000000)
	clientConfig.SetMaxWorkers(100)

	// Create the Unmeshed client
	var err error
	unmeshedClient, err = apis.NewUnmeshedClient(clientConfig)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// Register the single worker
	unmeshedClient.RegisterWorker(worker)

	// Create a channel to signal when the client is done
	done := make(chan struct{})

	// Start the client in a goroutine
	go func() {
		unmeshedClient.Start()
		close(done)
	}()

	fmt.Println("Single worker client started. Press Ctrl+C to stop...")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for either a shutdown signal or for the client to finish
	select {
	case <-sigChan:
		fmt.Println("\nReceived shutdown signal. Stopping client...")
	case <-done:
		fmt.Println("Client finished execution")
	}
}
