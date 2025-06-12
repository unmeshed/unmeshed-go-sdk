package common

type StepID struct {
	ID int64 `json:"id,omitempty"`
}

func NewStepID(id int64) *StepID {
	return &StepID{
		ID: id,
	}
}

func (s *StepID) GetID() int64 {
	return s.ID
}

func (s *StepID) SetID(id int64) {
	s.ID = id
}

type StepData struct {
	StepID        int64                  `json:"stepId,omitempty"`
	StepName      string                 `json:"stepName,omitempty"`
	StepRef       string                 `json:"stepRef,omitempty"`
	StepNamespace string                 `json:"stepNamespace,omitempty"`
	InputParam    map[string]interface{} `json:"inputParam,omitempty"`
	Output        map[string]interface{} `json:"output,omitempty"`
	Status        StepStatus             `json:"status,omitempty"`
	IsOptional    bool                   `json:"isOptional,omitempty"`
	Polling       int64                  `json:"polling,omitempty"`
	Scheduled     int64                  `json:"scheduled,omitempty"`
	Updated       int64                  `json:"updated,omitempty"`
	Priority      int64                  `json:"priority,omitempty"`
}

func NewStepData() *StepData {
	return &StepData{
		InputParam: make(map[string]interface{}),
		Output:     make(map[string]interface{}),
	}
}

// Getters
func (s *StepData) GetStepID() int64 {
	return s.StepID
}

func (s *StepData) GetStepName() string {
	return s.StepName
}

func (s *StepData) GetStepRef() string {
	return s.StepRef
}

func (s *StepData) GetStepNamespace() string {
	return s.StepNamespace
}

func (s *StepData) GetInputParam() map[string]interface{} {
	if s.InputParam == nil {
		s.InputParam = make(map[string]interface{})
	}
	return s.InputParam
}

func (s *StepData) GetOutput() map[string]interface{} {
	if s.Output == nil {
		s.Output = make(map[string]interface{})
	}
	return s.Output
}

func (s *StepData) GetStatus() StepStatus {
	return s.Status
}

func (s *StepData) GetIsOptional() bool {
	return s.IsOptional
}

func (s *StepData) GetPolling() int64 {
	return s.Polling
}

func (s *StepData) GetScheduled() int64 {
	return s.Scheduled
}

func (s *StepData) GetUpdated() int64 {
	return s.Updated
}

func (s *StepData) GetPriority() int64 {
	return s.Priority
}

// Setters
func (s *StepData) SetStepID(stepID int64) {
	s.StepID = stepID
}

func (s *StepData) SetStepName(stepName string) {
	s.StepName = stepName
}

func (s *StepData) SetStepRef(stepRef string) {
	s.StepRef = stepRef
}

func (s *StepData) SetStepNamespace(stepNamespace string) {
	s.StepNamespace = stepNamespace
}

func (s *StepData) SetInputParam(inputParam map[string]interface{}) {
	s.InputParam = inputParam
}

func (s *StepData) SetOutput(output map[string]interface{}) {
	s.Output = output
}

func (s *StepData) SetStatus(status StepStatus) {
	s.Status = status
}

func (s *StepData) SetIsOptional(isOptional bool) {
	s.IsOptional = isOptional
}

func (s *StepData) SetPolling(polling int64) {
	s.Polling = polling
}

func (s *StepData) SetScheduled(scheduled int64) {
	s.Scheduled = scheduled
}

func (s *StepData) SetUpdated(updated int64) {
	s.Updated = updated
}

func (s *StepData) SetPriority(priority int64) {
	s.Priority = priority
}
