package common

type WebhookSource string

const (
	WebhookSourceMSTeams    WebhookSource = "MS_TEAMS"
	WebhookSourceNotDefined WebhookSource = "NOT_DEFINED"
)

type SQAuthUserType string

const (
	SQAuthUserTypeUser     SQAuthUserType = "USER"
	SQAuthUserTypeAPI      SQAuthUserType = "API"
	SQAuthUserTypeInternal SQAuthUserType = "INTERNAL"
)

type StepType string

const (
	StepTypeWorker         StepType = "WORKER"
	StepTypeHTTP           StepType = "HTTP"
	StepTypeWait           StepType = "WAIT"
	StepTypeFail           StepType = "FAIL"
	StepTypePython         StepType = "PYTHON"
	StepTypeJavascript     StepType = "JAVASCRIPT"
	StepTypeJQ             StepType = "JQ"
	StepTypeManaged        StepType = "MANAGED"
	StepTypeBuiltin        StepType = "BUILTIN"
	StepTypeNoop           StepType = "NOOP"
	StepTypePersistedState StepType = "PERSISTED_STATE"
	StepTypeDependsOn      StepType = "DEPENDSON"
	StepTypeIntegration    StepType = "INTEGRATION"
	StepTypeExit           StepType = "EXIT"
	StepTypeSubProcess     StepType = "SUB_PROCESS"
	StepTypeList           StepType = "LIST"
	StepTypeParallel       StepType = "PARALLEL"
	StepTypeForeach        StepType = "FOREACH"
	StepTypeSwitch         StepType = "SWITCH"
	StepTypeScheduler      StepType = "SCHEDULER"
	StepTypeProcessTracker StepType = "PROCESS_TRACKER"
)

type StepStatus string

const (
	StepStatusPending   StepStatus = "PENDING"
	StepStatusScheduled StepStatus = "SCHEDULED"
	StepStatusRunning   StepStatus = "RUNNING"
	StepStatusPaused    StepStatus = "PAUSED"
	StepStatusCompleted StepStatus = "COMPLETED"
	StepStatusFailed    StepStatus = "FAILED"
	StepStatusTimedOut  StepStatus = "TIMED_OUT"
	StepStatusSkipped   StepStatus = "SKIPPED"
	StepStatusCancelled StepStatus = "CANCELLED"
)

type ProcessStatus string

const (
	ProcessStatusRunning    ProcessStatus = "RUNNING"
	ProcessStatusCompleted  ProcessStatus = "COMPLETED"
	ProcessStatusFailed     ProcessStatus = "FAILED"
	ProcessStatusTimedOut   ProcessStatus = "TIMED_OUT"
	ProcessStatusCancelled  ProcessStatus = "CANCELLED"
	ProcessStatusTerminated ProcessStatus = "TERMINATED"
	ProcessStatusReviewed   ProcessStatus = "REVIEWED"
)

type ProcessTriggerType string

const (
	ProcessTriggerTypeManual     ProcessTriggerType = "MANUAL"
	ProcessTriggerTypeScheduled  ProcessTriggerType = "SCHEDULED"
	ProcessTriggerTypeAPIMapping ProcessTriggerType = "API_MAPPING"
	ProcessTriggerTypeWebhook    ProcessTriggerType = "WEBHOOK"
	ProcessTriggerTypeAPI        ProcessTriggerType = "API"
	ProcessTriggerTypeSubProcess ProcessTriggerType = "SUB_PROCESS"
)

type ProcessType string

const (
	ProcessTypeStandard         ProcessType = "STANDARD"
	ProcessTypeDynamic          ProcessType = "DYNAMIC"
	ProcessTypeAPIOrchestration ProcessType = "API_ORCHESTRATION"
	ProcessTypeInternal         ProcessType = "INTERNAL"
)
