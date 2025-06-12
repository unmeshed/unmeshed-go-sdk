package common

type ProcessRequestData struct {
	Namespace     *string                `json:"namespace,omitempty"`
	Name          *string                `json:"name,omitempty"`
	Version       *int                   `json:"version,omitempty"`
	RequestID     *string                `json:"requestId,omitempty"`
	CorrelationID *string                `json:"correlationId,omitempty"`
	Input         map[string]interface{} `json:"input,omitempty"`
}

func NewProcessRequestData() *ProcessRequestData {
	return &ProcessRequestData{
		Input: make(map[string]interface{}),
	}
}

func (p *ProcessRequestData) GetNamespace() *string {
	return p.Namespace
}

func (p *ProcessRequestData) GetName() *string {
	return p.Name
}

func (p *ProcessRequestData) GetVersion() *int {
	return p.Version
}

func (p *ProcessRequestData) GetRequestID() *string {
	return p.RequestID
}

func (p *ProcessRequestData) GetCorrelationID() *string {
	return p.CorrelationID
}

func (p *ProcessRequestData) GetInput() map[string]interface{} {
	if p.Input == nil {
		p.Input = make(map[string]interface{})
	}
	return p.Input
}

func (p *ProcessRequestData) SetNamespace(namespace *string) {
	p.Namespace = namespace
}

func (p *ProcessRequestData) SetName(name *string) {
	p.Name = name
}

func (p *ProcessRequestData) SetVersion(version *int) {
	p.Version = version
}

func (p *ProcessRequestData) SetRequestID(requestID *string) {
	p.RequestID = requestID
}

func (p *ProcessRequestData) SetCorrelationID(correlationID *string) {
	p.CorrelationID = correlationID
}

func (p *ProcessRequestData) SetInput(input map[string]interface{}) {
	p.Input = input
}
