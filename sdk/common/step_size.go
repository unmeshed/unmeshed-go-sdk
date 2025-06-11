package common

type StepSize struct {
	StepQueueNameData StepQueueNameData `json:"stepQueueNameData,omitempty"`
	Size              int               `json:"size,omitempty"`
}

func NewStepSize(stepQueueNameData StepQueueNameData, size int) StepSize {
	return StepSize{
		StepQueueNameData: stepQueueNameData,
		Size:              size,
	}
}
