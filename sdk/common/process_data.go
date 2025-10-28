package common

type ProcessData struct {
	ProcessID       int64                  `json:"processId,omitempty"`
	ProcessType     *ProcessType           `json:"processType,omitempty"`
	TriggerType     *ProcessTriggerType    `json:"triggerType,omitempty"`
	Namespace       *string                `json:"namespace,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Version         *int                   `json:"version,omitempty"`
	HistoryID       *int                   `json:"historyId,omitempty"`
	RequestID       *string                `json:"requestId,omitempty"`
	CorrelationID   *string                `json:"correlationId,omitempty"`
	Status          *ProcessStatus         `json:"status,omitempty"`
	Input           map[string]interface{} `json:"input,omitempty"`
	Output          map[string]interface{} `json:"output,omitempty"`
	State           map[string]interface{} `json:"state,omitempty"`
	SecretState     map[string]interface{} `json:"secretState,omitempty"`
	AuthClaims      map[string]interface{} `json:"authClaims,omitempty"`
	StepIDCount     *int                   `json:"stepIdCount,omitempty"`
	ShardName       *string                `json:"shardName,omitempty"`
	ShardInstanceID *int                   `json:"shardInstanceId,omitempty"`
	Steps           []*StepID              `json:"steps,omitempty"`
	StepRecords     []*StepData            `json:"stepRecords,omitempty"`
	Created         *int64                 `json:"created,omitempty"`
	Updated         *int64                 `json:"updated,omitempty"`
	CreatedBy       *string                `json:"createdBy,omitempty"`
	Tags            []*TagValue            `json:"tags,omitempty"`
}

func NewProcessData() *ProcessData {
	return &ProcessData{
		ProcessID: 0,
	}
}

// Getters
func (p *ProcessData) GetProcessID() int64 {
	return p.ProcessID
}

func (p *ProcessData) GetProcessType() *ProcessType {
	return p.ProcessType
}

func (p *ProcessData) GetTriggerType() *ProcessTriggerType {
	return p.TriggerType
}

func (p *ProcessData) GetNamespace() *string {
	return p.Namespace
}

func (p *ProcessData) GetName() *string {
	return p.Name
}

func (p *ProcessData) GetVersion() *int {
	return p.Version
}

func (p *ProcessData) GetHistoryID() *int {
	return p.HistoryID
}

func (p *ProcessData) GetRequestID() *string {
	return p.RequestID
}

func (p *ProcessData) GetCorrelationID() *string {
	return p.CorrelationID
}

func (p *ProcessData) GetStatus() *ProcessStatus {
	return p.Status
}

func (p *ProcessData) GetInput() map[string]interface{} {
	return p.Input
}

func (p *ProcessData) GetOutput() map[string]interface{} {
	return p.Output
}

func (p *ProcessData) GetState() map[string]interface{} {
	return p.State
}

func (p *ProcessData) GetSecretState() map[string]interface{} {
	return p.SecretState
}

func (p *ProcessData) GetAuthClaims() map[string]interface{} {
	return p.AuthClaims
}

func (p *ProcessData) GetStepIDCount() *int {
	return p.StepIDCount
}

func (p *ProcessData) GetShardName() *string {
	return p.ShardName
}

func (p *ProcessData) GetShardInstanceID() *int {
	return p.ShardInstanceID
}

func (p *ProcessData) GetSteps() []*StepID {
	return p.Steps
}

func (p *ProcessData) GetStepRecords() []*StepData {
	return p.StepRecords
}

func (p *ProcessData) GetCreated() *int64 {
	return p.Created
}

func (p *ProcessData) GetUpdated() *int64 {
	return p.Updated
}

func (p *ProcessData) GetCreatedBy() *string {
	return p.CreatedBy
}

// Setters
func (p *ProcessData) SetProcessID(processID int64) {
	p.ProcessID = processID
}

func (p *ProcessData) SetProcessType(processType *ProcessType) {
	p.ProcessType = processType
}

func (p *ProcessData) SetTriggerType(triggerType *ProcessTriggerType) {
	p.TriggerType = triggerType
}

func (p *ProcessData) SetNamespace(namespace *string) {
	p.Namespace = namespace
}

func (p *ProcessData) SetName(name *string) {
	p.Name = name
}

func (p *ProcessData) SetVersion(version *int) {
	p.Version = version
}

func (p *ProcessData) SetHistoryID(historyID *int) {
	p.HistoryID = historyID
}

func (p *ProcessData) SetRequestID(requestID *string) {
	p.RequestID = requestID
}

func (p *ProcessData) SetCorrelationID(correlationID *string) {
	p.CorrelationID = correlationID
}

func (p *ProcessData) SetStatus(status *ProcessStatus) {
	p.Status = status
}

func (p *ProcessData) SetInput(input map[string]interface{}) {
	p.Input = input
}

func (p *ProcessData) SetOutput(output map[string]interface{}) {
	p.Output = output
}

func (p *ProcessData) SetState(state map[string]interface{}) {
	p.State = state
}

func (p *ProcessData) SetSecretState(secretState map[string]interface{}) {
	p.SecretState = secretState
}

func (p *ProcessData) SetAuthClaims(authClaims map[string]interface{}) {
	p.AuthClaims = authClaims
}

func (p *ProcessData) SetStepIDCount(stepIDCount *int) {
	p.StepIDCount = stepIDCount
}

func (p *ProcessData) SetShardName(shardName *string) {
	p.ShardName = shardName
}

func (p *ProcessData) SetShardInstanceID(shardInstanceID *int) {
	p.ShardInstanceID = shardInstanceID
}

func (p *ProcessData) SetSteps(steps []*StepID) {
	p.Steps = steps
}

func (p *ProcessData) SetStepRecords(stepRecords []*StepData) {
	p.StepRecords = stepRecords
}

func (p *ProcessData) SetCreated(created *int64) {
	p.Created = created
}

func (p *ProcessData) SetUpdated(updated *int64) {
	p.Updated = updated
}

func (p *ProcessData) SetCreatedBy(createdBy *string) {
	p.CreatedBy = createdBy
}

func (p *ProcessData) GetTags() []*TagValue {
	return p.Tags
}

// Setter for Tags
func (p *ProcessData) SetTags(tags []*TagValue) {
	p.Tags = tags
}
