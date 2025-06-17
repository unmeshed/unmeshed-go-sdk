package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/main"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

func TestNewUnmeshedClient(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")

	client, err := apis.NewUnmeshedClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, config, client.ClientConfig)
	assert.NotNil(t, client.DoneChan())
}

func TestNewUnmeshedClient_InvalidConfig(t *testing.T) {
	config := &configs.ClientConfig{}

	client, err := apis.NewUnmeshedClient(config)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "cannot initialize without a valid clientId and token")

	config.SetClientID("test-client")
	client, err = apis.NewUnmeshedClient(config)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "cannot initialize without a valid clientId and token")
}

func TestRegisterWorker_SuccessAndDuplicate(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client, err := apis.NewUnmeshedClient(config)
	assert.NoError(t, err)

	worker := workers.NewWorker(func(input interface{}) interface{} { return input }, "worker1")
	err = client.RegisterWorker(worker)
	assert.NoError(t, err)

	err = client.RegisterWorker(worker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestRegisterWorkers_Multiple(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client, err := apis.NewUnmeshedClient(config)
	assert.NoError(t, err)

	worker1 := workers.NewWorker(func(input interface{}) interface{} { return input }, "worker1")
	worker2 := workers.NewWorker(func(input interface{}) interface{} { return input }, "worker2")
	workersList := []*workers.Worker{worker1, worker2}

	err = client.RegisterWorkers(workersList)
	assert.NoError(t, err)

	err = client.RegisterWorker(worker1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestRegisterWorker_NoExecutionMethod(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client, err := apis.NewUnmeshedClient(config)
	assert.NoError(t, err)

	worker := workers.NewWorker(nil, "worker1")
	err = client.RegisterWorker(worker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no execution method found")
}

func TestRegisterWorker_WrongSignature(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client, err := apis.NewUnmeshedClient(config)
	assert.NoError(t, err)

	worker := workers.NewWorker(func(a, b interface{}) interface{} { return a }, "worker1")
	err = client.RegisterWorker(worker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must have exactly one parameter")
}

func TestStart_NoWorkers(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client, err := apis.NewUnmeshedClient(config)
	assert.NoError(t, err)

	ch := make(chan struct{})
	go func() {
		client.Start()
		close(ch)
	}()
	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		t.Fatal("Start did not return when no workers are configured")
	}
}

func TestStop_MultipleCalls(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client, err := apis.NewUnmeshedClient(config)
	assert.NoError(t, err)

	// First call to Stop
	client.Stop()

	// Verify done channel is closed
	select {
	case <-client.DoneChan():
	default:
		t.Error("done channel should be closed after Stop()")
	}

	client.Stop()
}

func TestProcessOperations(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client, err := apis.NewUnmeshedClient(config)
	assert.NoError(t, err)

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
		},
	}

	_, err = client.RunProcessSyncWithDefaultTimeout(processRequest)
	assert.NotNil(t, err)

	_, err = client.RunProcessAsync(processRequest)
	assert.NotNil(t, err)

	searchParams := &common.ProcessSearchRequest{
		Names:     []string{"test_process"},
		Limit:     20,
		Namespace: &namespace,
	}
	_, err = client.SearchProcessExecutions(searchParams)
	assert.NotNil(t, err)
}
