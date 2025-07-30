package common

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

type WorkResponseBuilder struct{}

func NewWorkResponseBuilder() *WorkResponseBuilder {
	return &WorkResponseBuilder{}
}

func (b *WorkResponseBuilder) resultToMap(obj interface{}) map[string]interface{} {
	if obj == nil {
		return map[string]interface{}{"result": nil}
	}

	switch v := obj.(type) {
	case bool, string, int, int64, float64:
		return map[string]interface{}{"result": v}

	case map[string]interface{}:
		return v
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			strKey, ok := key.(string)
			if !ok {
				log.Printf("Non-string key found in map: %v", key)
				continue
			}
			result[strKey] = value
		}
		return result

	case []interface{}:
		return map[string]interface{}{"result": v}
	case []string:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = value
		}
		return map[string]interface{}{"result": result}
	case []int:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = value
		}
		return map[string]interface{}{"result": result}
	case []int64:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = value
		}
		return map[string]interface{}{"result": result}
	case []float64:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = value
		}
		return map[string]interface{}{"result": result}

	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			log.Printf("Error marshaling object: %v", err)
			return map[string]interface{}{"error": "failed to marshal"}
		}

		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			return map[string]interface{}{"error": "failed to unmarshal"}
		}

		return result
	}
}

func (b *WorkResponseBuilder) FailResponse(workRequest *WorkRequest, context error) *WorkResponse {
	actualCause := b.tryPeelIrrelevantExceptions(context)

	var innerError interface{}
	errMsg := actualCause.Error()

	if json.Unmarshal([]byte(errMsg), &innerError) != nil {
		innerError = errMsg
	}

	output := map[string]interface{}{
		"error": innerError,
	}

	workResponse := NewWorkResponse()
	workResponse.SetProcessID(workRequest.ProcessID)
	workResponse.SetStepID(workRequest.GetStepID())
	if workRequest.StepExecutionID == 0 {
		workResponse.StepExecutionID = 0
	} else {
		workResponse.SetStepExecutionID(workRequest.GetStepExecutionID())
	}
	workResponse.SetOutput(output)
	workResponse.SetStartedAt(time.Now().UnixMilli())
	workResponse.SetStatus(StepStatusFailed)

	return workResponse
}

func (b *WorkResponseBuilder) tryPeelIrrelevantExceptions(context error) error {
	actualCause := context

	if future, ok := context.(interface{ Load() any }); ok {
		if val := future.Load(); val != nil {
			if err, ok := val.(error); ok {
				unwrapped := errors.Unwrap(err)
				if unwrapped != nil {
					return unwrapped
				}
				return err
			}
		}
	}

	return actualCause
}

func (b *WorkResponseBuilder) SuccessResponse(workRequest *WorkRequest, stepResult *StepResult) *WorkResponse {
	output := b.resultToMap(stepResult.GetResult())

	workResponse := NewWorkResponse()
	workResponse.SetProcessID(workRequest.GetProcessID())
	workResponse.SetStepID(workRequest.GetStepID())
	if workRequest.StepExecutionID == 0 {
		workResponse.StepExecutionID = 0
	} else {
		workResponse.SetStepExecutionID(workRequest.GetStepExecutionID())
	}
	workResponse.SetOutput(output)
	workResponse.SetStartedAt(time.Now().UnixMilli())
	workResponse.SetStatus(StepStatusCompleted)

	return workResponse
}

func (b *WorkResponseBuilder) RunningResponse(workRequest *WorkRequest, stepResult *StepResult) *WorkResponse {
	output := b.resultToMap(stepResult.GetResult())
	workResponse := NewWorkResponse()
	workResponse.SetProcessID(workRequest.GetProcessID())
	workResponse.SetStepID(workRequest.GetStepID())
	if workRequest.StepExecutionID == 0 {
		workResponse.StepExecutionID = 0
	} else {
		workResponse.SetStepExecutionID(workRequest.GetStepExecutionID())
	}
	workResponse.SetOutput(output)
	workResponse.SetStartedAt(time.Now().UnixMilli())
	workResponse.SetStatus(StepStatusRunning)
	workResponse.SetRescheduleAfterSeconds(stepResult.RescheduleAfterSeconds)
	return workResponse
}
