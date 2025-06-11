package common

type WorkRequest struct {
	ProcessID       int64                  `json:"processId,omitempty"`
	StepID          int64                  `json:"stepId,omitempty"`
	StepExecutionID int64                  `json:"stepExecutionId"`
	StepName        string                 `json:"stepName,omitempty"`
	StepRef         string                 `json:"stepRef,omitempty"`
	StepNamespace   string                 `json:"stepNamespace,omitempty"`
	InputParam      map[string]interface{} `json:"inputParam,omitempty"`
	IsOptional      bool                   `json:"isOptional,omitempty"`
	Polling         int64                  `json:"polling,omitempty"`
	Scheduled       int64                  `json:"scheduled,omitempty"`
	Updated         int64                  `json:"updated,omitempty"`
	Priority        int64                  `json:"priority,omitempty"`
}

func NewWorkRequest() *WorkRequest {
	return &WorkRequest{
		InputParam: make(map[string]interface{}),
	}
}

func (w *WorkRequest) SetStepNamespace(stepNamespace string) {
	w.StepNamespace = stepNamespace
}

func (w *WorkRequest) GetStepNamespace() string {
	return w.StepNamespace
}

func (w *WorkRequest) GetStepName() string {
	return w.StepName
}

func (w *WorkRequest) GetProcessID() int64 {
	return w.ProcessID
}

func (w *WorkRequest) GetStepID() int64 {
	return w.StepID
}

func (w *WorkRequest) GetStepExecutionID() int64 {
	return w.StepExecutionID
}

func (w *WorkRequest) GetPolling() int64 {
	return w.Polling
}

func (w *WorkRequest) GetScheduled() int64 {
	return w.Scheduled
}

func (w *WorkRequest) GetUpdated() int64 {
	return w.Updated
}

func (w *WorkRequest) GetPriority() int64 {
	return w.Priority
}

func (w *WorkRequest) GetInputParam() map[string]interface{} {
	if w.InputParam == nil {
		w.InputParam = make(map[string]interface{})
	}
	return w.InputParam
}
