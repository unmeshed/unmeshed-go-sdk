package common

type StepResult struct {
	Result                 interface{}
	KeepRunning            bool
	RescheduleAfterSeconds int
}

func NewStepResult(result interface{}) *StepResult {
	return &StepResult{Result: result, KeepRunning: false, RescheduleAfterSeconds: 0}
}

func (sr *StepResult) GetResult() interface{} {
	return sr.Result
}
