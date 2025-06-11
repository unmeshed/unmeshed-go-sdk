package common

type StepQueueNameData struct {
	OrgId     int      `json:"orgId,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
	StepType  StepType `json:"stepType,omitempty"`
	Name      string   `json:"name,omitempty"`
}

func NewStepQueueNameData(orgId int, namespace, name string, stepType StepType) *StepQueueNameData {
	return &StepQueueNameData{
		OrgId:     orgId,
		Namespace: namespace,
		Name:      name,
		StepType:  stepType,
	}
}
