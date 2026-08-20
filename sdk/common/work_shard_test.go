package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkRequestDecodesShardInstanceID(t *testing.T) {
	var request WorkRequest

	require.NoError(t, json.Unmarshal([]byte(`{"stepId":1,"stepExecutionId":2,"shardInstanceId":7}`), &request))

	require.NotNil(t, request.GetShardInstanceID())
	assert.Equal(t, 7, *request.GetShardInstanceID())
}

func TestWorkResponseBuilderCopiesShardInstanceID(t *testing.T) {
	shardInstanceID := 7
	request := NewWorkRequest()
	request.ProcessID = 11
	request.StepID = 22
	request.StepExecutionID = 33
	request.SetShardInstanceID(&shardInstanceID)

	response := NewWorkResponseBuilder().SuccessResponse(request, NewStepResult("ok"))

	require.NotNil(t, response.GetShardInstanceID())
	assert.Equal(t, 7, *response.GetShardInstanceID())
}
