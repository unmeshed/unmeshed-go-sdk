package common

import (
	"fmt"
)

type PollRequestData struct {
	StepQueueNameData *StepQueueNameData
	Size              int64
}

func NewPollRequestData(stepQueueNameData *StepQueueNameData, size int64) (*PollRequestData, error) {
	prd := &PollRequestData{}

	if stepQueueNameData == nil {
		return nil, fmt.Errorf("stepQueueNameData cannot be nil")
	}
	prd.StepQueueNameData = stepQueueNameData

	if err := prd.SetSize(size); err != nil {
		return nil, err
	}

	return prd, nil
}

func (prd *PollRequestData) GetStepQueueNameData() *StepQueueNameData {
	return prd.StepQueueNameData
}

// SetStepQueueNameData sets the StepQueueNameData.
func (prd *PollRequestData) SetStepQueueNameData(value *StepQueueNameData) error {
	if value == nil {
		return fmt.Errorf("stepQueueNameData cannot be nil")
	}
	prd.StepQueueNameData = value
	return nil
}

func (prd *PollRequestData) GetSize() int64 {
	return prd.Size
}

func (prd *PollRequestData) SetSize(value int64) error {
	if value < 0 {
		return fmt.Errorf("size cannot be negative")
	}
	prd.Size = value
	return nil
}

func (prd *PollRequestData) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"stepQueueNameData": prd.StepQueueNameData,
		"size":              prd.Size,
	}
}
