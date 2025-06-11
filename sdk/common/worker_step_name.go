package common

type WorkerStepName struct {
	StepQueueNameData StepQueueNameData `json:"stepQueueNameData,omitempty"`
}

func NewWorkerStepName(StepQueueNameData StepQueueNameData) *WorkerStepName {
	return &WorkerStepName{
		StepQueueNameData: StepQueueNameData,
	}
}

func (wsn *WorkerStepName) GetStepQueueNameData() StepQueueNameData {
	return wsn.StepQueueNameData
}

func (wsn *WorkerStepName) SetStepQueueNameData(stepQueueNameData StepQueueNameData) {
	wsn.StepQueueNameData = stepQueueNameData
}
