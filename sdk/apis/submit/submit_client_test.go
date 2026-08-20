package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apisHttp "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/http"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

type submitRequest struct {
	shardHeader string
	stepIDs     []int64
}

func TestProcessBatchByShardUsesShardHeader(t *testing.T) {
	t.Setenv("DISABLE_SUBMIT_CLIENT", "true")

	var mu sync.Mutex
	requests := make([]submitRequest, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/clients/bulkResults", r.URL.Path)

		var batch []common.WorkResponse
		require.NoError(t, json.NewDecoder(r.Body).Decode(&batch))

		stepIDs := make([]int64, 0, len(batch))
		results := make(map[string]*common.ClientSubmitResult, len(batch))
		for _, response := range batch {
			stepIDs = append(stepIDs, response.StepID)
			results[fmt.Sprintf("%d", response.StepID)] = common.NewClientSubmitResult(response.ProcessID, response.StepID, http.StatusOK, "")
		}

		mu.Lock()
		requests = append(requests, submitRequest{
			shardHeader: r.Header.Get(SHARD_ID_HEADER),
			stepIDs:     stepIDs,
		})
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(w).Encode(results))
	}))
	defer server.Close()

	config := configs.NewClientConfig()
	config.SetClientID("test-client")
	config.SetBaseURL(server.URL)
	client := NewSubmitClient(apisHttp.NewHttpRequestFactory(config), config)
	defer client.Stop()

	shardSeven := 7
	client.processBatchByShard([]*common.WorkResponse{
		workResponseForTest(1, &shardSeven),
		workResponseForTest(2, nil),
		workResponseForTest(3, &shardSeven),
	})

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, requests, 2)

	byHeader := map[string][]int64{}
	for _, request := range requests {
		byHeader[request.shardHeader] = request.stepIDs
	}

	assert.ElementsMatch(t, []int64{1, 3}, byHeader["shard-7"])
	assert.ElementsMatch(t, []int64{2}, byHeader[""])
}

func workResponseForTest(stepID int64, shardInstanceID *int) *common.WorkResponse {
	response := common.NewWorkResponse()
	response.SetProcessID(100 + stepID)
	response.SetStepID(stepID)
	response.SetStepExecutionID(1000 + stepID)
	response.SetShardInstanceID(shardInstanceID)
	response.SetStatus(common.StepStatusCompleted)
	return response
}
