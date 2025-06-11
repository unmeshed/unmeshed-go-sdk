package common

import "encoding/json"

type WorkResponse struct {
	ProcessID              int64                  `json:"processId,omitempty"`
	StepID                 int64                  `json:"stepId,omitempty"`
	StepExecutionID        int64                  `json:"stepExecutionId"`
	Output                 map[string]interface{} `json:"output,omitempty"`
	Status                 StepStatus             `json:"status,omitempty"`
	RescheduleAfterSeconds int                    `json:"rescheduleAfterSeconds"`
	StartedAt              int64                  `json:"startedAt,omitempty"`
}

func NewWorkResponse() *WorkResponse {
	return &WorkResponse{
		Output: make(map[string]interface{}),
	}
}

func (wr *WorkResponse) GetProcessID() int64 {
	return wr.ProcessID
}

func (wr *WorkResponse) SetProcessID(processID int64) {
	wr.ProcessID = processID
}

func (wr *WorkResponse) GetStepID() int64 {
	return wr.StepID
}

func (wr *WorkResponse) SetStepID(stepID int64) {
	wr.StepID = stepID
}

func (wr *WorkResponse) GetStepExecutionID() int64 {
	return wr.StepExecutionID
}

func (wr *WorkResponse) SetStepExecutionID(stepExecutionID int64) {
	wr.StepExecutionID = stepExecutionID
}

func (wr *WorkResponse) GetOutput() map[string]interface{} {
	return wr.Output
}

func (wr *WorkResponse) SetOutput(output map[string]interface{}) {
	wr.Output = output
}

func (wr *WorkResponse) GetStatus() StepStatus {
	return wr.Status
}

func (wr *WorkResponse) SetStatus(status StepStatus) {
	wr.Status = status
}

func (wr *WorkResponse) GetRescheduleAfterSeconds() int {
	return wr.RescheduleAfterSeconds
}

func (wr *WorkResponse) SetRescheduleAfterSeconds(rescheduleAfterSeconds int) {
	wr.RescheduleAfterSeconds = rescheduleAfterSeconds
}

func (wr *WorkResponse) GetStartedAt() int64 {
	return wr.StartedAt
}

func (wr *WorkResponse) SetStartedAt(startedAt int64) {
	wr.StartedAt = startedAt
}

func (wr *WorkResponse) ToJSON() (string, error) {
	data, err := json.Marshal(wr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
