package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
)

type FunctionWrapper struct {
	Fn  interface{}
	Arg interface{}
}

type WorkerRunner struct{}

func NewWorkerRunner() *WorkerRunner {
	return &WorkerRunner{}
}

func (wr *WorkerRunner) RunWorker(worker *workers.Worker, workRequest *common.WorkRequest) (interface{}, error) {
	wrapper := FunctionWrapper{
		Fn:  worker.ExecutionMethod,
		Arg: workRequest.InputParam,
	}
	result, err := wr.invokeFunction(wrapper)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (wr *WorkerRunner) invokeFunction(f FunctionWrapper) (interface{}, error) {
	fnType := reflect.TypeOf(f.Fn)

	if fnType.Kind() != reflect.Func {
		log.Printf("Skipping invalid function: %+v (Incorrect signature)\n", f.Fn)
		return nil, errors.New("Skipping invalid function")
	}

	argType := fnType.In(0)
	argValue := reflect.ValueOf(f.Arg)

	if argValue.Type().Kind() == reflect.Map && argValue.Type().Key().Kind() == reflect.String {
		targetValue := reflect.New(argType).Elem()
		jsonData, err := json.Marshal(f.Arg)
		if err != nil {
			log.Printf("JSON marshal error for function: %+v, error: %v\n", f.Fn, err)
			return nil, err
		}

		err = json.Unmarshal(jsonData, targetValue.Addr().Interface())
		if err != nil {
			log.Printf("JSON unmarshal error for function: %+v, error: %v\n", f.Fn, err)
			return nil, err
		}
		argValue = targetValue
	} else if argValue.Type().Kind() == reflect.Slice {
		targetValue := reflect.New(argType).Elem()
		jsonData, err := json.Marshal(f.Arg)
		if err != nil {
			log.Printf("JSON marshal error for function: %+v, error: %v\n", f.Fn, err)
			return nil, err
		}

		err = json.Unmarshal(jsonData, targetValue.Addr().Interface())
		if err != nil {
			log.Printf("JSON unmarshal error for function: %+v, error: %v\n", f.Fn, err)
			return nil, err
		}
		argValue = targetValue
	} else {
		log.Printf("Argument for function: %+v must be a map[string]interface{} or []interface{}\n", f.Fn)
		return nil, fmt.Errorf("Argument for function: %+v must be a map[string]interface{} or []interface{}", f.Fn)
	}

	// Ensure pointer compatibility
	if argType.Kind() == reflect.Ptr && argValue.Kind() != reflect.Ptr {
		argValue = argValue.Addr()
	} else if argType.Kind() != reflect.Ptr && argValue.Kind() == reflect.Ptr {
		argValue = argValue.Elem()
	}

	results := reflect.ValueOf(f.Fn).Call([]reflect.Value{argValue})

	var finalResults []interface{}

	for _, result := range results {
		finalResults = append(finalResults, result.Interface())
	}

	if len(finalResults) > 0 {
		if err, ok := finalResults[len(finalResults)-1].(error); ok {
			results := finalResults[:len(finalResults)-1]
			if err != nil {
				return nil, err
			}
			if len(results) == 1 {
				result := results[0]
				if reflect.TypeOf(result).Kind() != reflect.Slice && reflect.TypeOf(result).Kind() != reflect.Array {
					return result, err
				}
			}
			return results, err
		}
	}

	if len(finalResults) == 1 {
		result := finalResults[0]
		if reflect.TypeOf(result).Kind() != reflect.Slice && reflect.TypeOf(result).Kind() != reflect.Array {
			return result, nil
		}
	}

	return finalResults, nil
}

func (wr *WorkerRunner) invokeFunctions(functions []FunctionWrapper) {
	for _, f := range functions {
		wr.invokeFunction(f)
	}
}
