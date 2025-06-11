package common

type StepQueueName struct {
	OrgId     int      `json:"orgId"`
	Namespace string   `json:"namespace"`
	StepType  StepType `json:"stepType"`
	Name      string   `json:"name"`
}

func NewStepQueueName(orgId int, namespace string, stepType StepType, name string) *StepQueueName {
	return &StepQueueName{
		OrgId:     orgId,
		Namespace: namespace,
		StepType:  stepType,
		Name:      name,
	}
}
