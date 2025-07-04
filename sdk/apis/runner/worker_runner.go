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

	if fnType.NumIn() != 1 {
		log.Printf("Function must accept exactly one argument: %+v\n", f.Fn)
		return nil, errors.New("Function must accept exactly one argument")
	}

	argType := fnType.In(0)
	argValue := reflect.ValueOf(f.Arg)

	if argValue.Type().Kind() == reflect.Map && argValue.Type().Key().Kind() == reflect.String ||
		argValue.Type().Kind() == reflect.Slice {
		targetValue := reflect.New(argType).Elem()
		jsonData, err := json.Marshal(f.Arg)
		if err != nil {
			log.Printf("JSON marshal error: %v\n", err)
			return nil, err
		}
		err = json.Unmarshal(jsonData, targetValue.Addr().Interface())
		if err != nil {
			log.Printf("JSON unmarshal error: %v\n", err)
			return nil, err
		}
		argValue = targetValue
	} else {
		log.Printf("Invalid input type for function: %+v\n", f.Fn)
		return nil, fmt.Errorf("Argument must be map[string]interface{} or []interface{}")
	}

	if argType.Kind() == reflect.Ptr && argValue.Kind() != reflect.Ptr {
		argValue = argValue.Addr()
	} else if argType.Kind() != reflect.Ptr && argValue.Kind() == reflect.Ptr {
		argValue = argValue.Elem()
	}

	rawResults := reflect.ValueOf(f.Fn).Call([]reflect.Value{argValue})
	numResults := len(rawResults)

	if numResults == 0 {
		return nil, nil
	}

	// Detect if last value is an error (even if interface is nil)
	lastVal := rawResults[numResults-1]
	lastType := fnType.Out(numResults - 1)

	isError := lastType.Implements(reflect.TypeOf((*error)(nil)).Elem())
	finalResults := []interface{}{}

	if isError {
		errVal := lastVal.Interface()
		if errVal != nil {
			return nil, errVal.(error)
		}
		// Discard nil error
		rawResults = rawResults[:numResults-1]
		numResults--
	}

	// Build result slice
	for i := 0; i < numResults; i++ {
		finalResults = append(finalResults, rawResults[i].Interface())
	}

	if len(finalResults) == 1 {
		// Unwrap single result (not slice)
		kind := reflect.TypeOf(finalResults[0]).Kind()
		if kind != reflect.Slice && kind != reflect.Array {
			return finalResults[0], nil
		}
	}

	return finalResults, nil
}

func (wr *WorkerRunner) invokeFunctions(functions []FunctionWrapper) {
	for _, f := range functions {
		_, _ = wr.invokeFunction(f)
	}
}
