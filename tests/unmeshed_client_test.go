package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/main"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

func TestNewUnmeshedClient(t *testing.T) {
	config := &configs.ClientConfig{}

	client := apis.NewUnmeshedClient(config)
	assert.NotNil(t, client)
	assert.Equal(t, config, client.ClientConfig)
}

func TestNewUnmeshedClient_ValidConfig(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client := apis.NewUnmeshedClient(config)
	assert.NotNil(t, client)
	assert.Equal(t, config, client.ClientConfig)
}

func TestRegisterWorker_SuccessAndDuplicate(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client := apis.NewUnmeshedClient(config)

	worker := workers.NewWorker(func(input interface{}) interface{} { return input }, "worker1")
	err := client.RegisterWorker(worker)
	assert.NoError(t, err)

	// Duplicate registration should fail
	err = client.RegisterWorker(worker)
	assert.Error(t, err)
}

func TestRegisterWorkers_Multiple(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client := apis.NewUnmeshedClient(config)

	worker1 := workers.NewWorker(func(input interface{}) interface{} { return input }, "worker1")
	worker2 := workers.NewWorker(func(input interface{}) interface{} { return input }, "worker2")
	workersList := []*workers.Worker{worker1, worker2}

	err := client.RegisterWorkers(workersList)
	assert.NoError(t, err)
}

func TestRegisterWorker_NoExecutionMethod(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client := apis.NewUnmeshedClient(config)

	worker := workers.NewWorker(nil, "worker1")
	err := client.RegisterWorker(worker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no execution method found")
}

func TestRegisterWorker_WrongSignature(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client := apis.NewUnmeshedClient(config)

	worker := workers.NewWorker(func(a, b interface{}) interface{} { return a }, "worker1")
	err := client.RegisterWorker(worker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must have exactly one parameter")
}

func TestStart_NoWorkers(t *testing.T) {
	config := &configs.ClientConfig{}
	config.SetClientID("test-client")
	config.SetAuthToken("test-token")
	client := apis.NewUnmeshedClient(config)

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
	client := apis.NewUnmeshedClient(config)

	client.Stop()

	select {
	case <-client.DoneChan():
	default:
		t.Error("done channel should be closed after Stop()")
	}

	client.Stop()
}
