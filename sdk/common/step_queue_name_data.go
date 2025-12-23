package common

type StepQueueNameData struct {
	OrgId      int      `json:"orgId,omitempty"`
	Namespace  string   `json:"namespace,omitempty"`
	StepType   StepType `json:"stepType,omitempty"`
	Name       string   `json:"name,omitempty"`
	ServerName string   `json:"serverName,omitempty"`
}

func NewStepQueueNameData(orgId int, namespace, name string, stepType StepType, serverName string) *StepQueueNameData {
	return &StepQueueNameData{
		OrgId:      orgId,
		Namespace:  namespace,
		Name:       name,
		StepType:   stepType,
		ServerName: serverName,
	}
}
