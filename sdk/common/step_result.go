package common

type StepResult struct {
	result interface{}
}

func NewStepResult(result interface{}) *StepResult {
	return &StepResult{result: result}
}

func (sr *StepResult) GetResult() interface{} {
	return sr.result
}
