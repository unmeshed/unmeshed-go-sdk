package tests

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	apisHttp "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/http"
	apisSubmit "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/submit"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

func TestNewSubmitClient_Disabled(t *testing.T) {
	os.Setenv("DISABLE_SUBMIT_CLIENT", "true")
	defer os.Unsetenv("DISABLE_SUBMIT_CLIENT")
	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	factory := apisHttp.NewHttpRequestFactory(config)
	client := apisSubmit.NewSubmitClient(factory, config)
	assert.NotNil(t, client)
}

func TestNewSubmitClient_Enabled(t *testing.T) {
	os.Unsetenv("DISABLE_SUBMIT_CLIENT")
	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	factory := apisHttp.NewHttpRequestFactory(config)
	client := apisSubmit.NewSubmitClient(factory, config)
	assert.NotNil(t, client)
	client.Stop()
}

func TestSubmit_AddsToTrackerAndQueue(t *testing.T) {
	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	factory := apisHttp.NewHttpRequestFactory(config)
	client := apisSubmit.NewSubmitClient(factory, config)
	workResponse := common.NewWorkResponse()
	workResponse.SetStepID(123)
	stepPollState := common.NewStepPollState(1)
	err := client.Submit(workResponse, stepPollState)
	assert.NoError(t, err)
	assert.Equal(t, 1, client.GetSubmitTrackerSize())
	client.Stop()
}

func TestSubmit_DuplicateStepID(t *testing.T) {
	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	factory := apisHttp.NewHttpRequestFactory(config)
	client := apisSubmit.NewSubmitClient(factory, config)
	workResponse1 := common.NewWorkResponse()
	workResponse1.SetStepID(123)
	stepPollState1 := common.NewStepPollState(1)
	client.Submit(workResponse1, stepPollState1)
	workResponse2 := common.NewWorkResponse()
	workResponse2.SetStepID(123)
	stepPollState2 := common.NewStepPollState(1)
	client.Submit(workResponse2, stepPollState2)
	assert.Equal(t, 1, client.GetSubmitTrackerSize())
	client.Stop()
}

func TestSubmit_StepPollStateCount(t *testing.T) {
	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	factory := apisHttp.NewHttpRequestFactory(config)
	client := apisSubmit.NewSubmitClient(factory, config)
	workResponse := common.NewWorkResponse()
	workResponse.SetStepID(123)
	stepPollState := common.NewStepPollState(3)
	err := client.Submit(workResponse, stepPollState)
	assert.NoError(t, err)
	assert.Equal(t, 1, client.GetSubmitTrackerSize())
	client.Stop()
}

func TestSubmit_EmptyClientIDPanics(t *testing.T) {
	if os.Getenv("TEST_EMPTY_CLIENT_ID") == "1" {
		config := configs.NewClientConfig()
		factory := apisHttp.NewHttpRequestFactory(config)
		_ = apisSubmit.NewSubmitClient(factory, config)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestSubmit_EmptyClientIDPanics")
	cmd.Env = append(os.Environ(), "TEST_EMPTY_CLIENT_ID=1")
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok && !exitErr.Success() {
		return
	}
	t.Fatalf("expected process to exit with non-zero status, got err=%v", err)
}

func TestGetSubmitTrackerSize(t *testing.T) {
	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	factory := apisHttp.NewHttpRequestFactory(config)
	client := apisSubmit.NewSubmitClient(factory, config)
	assert.Equal(t, 0, client.GetSubmitTrackerSize())
	workResponse := common.NewWorkResponse()
	workResponse.SetStepID(1)
	stepPollState := common.NewStepPollState(1)
	client.Submit(workResponse, stepPollState)
	assert.Equal(t, 1, client.GetSubmitTrackerSize())
	client.Stop()
}

func TestClose_Idempotent(t *testing.T) {
	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	factory := apisHttp.NewHttpRequestFactory(config)
	client := apisSubmit.NewSubmitClient(factory, config)
	client.Stop()
	client.Stop()
}
