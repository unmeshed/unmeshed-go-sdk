package common

type ProcessConfiguration struct {
	CompletionTimeout int                 `json:"completionTimeout" default:"180000"`
	OnTimeoutProcess  *ProcessRequestData `json:"onTimeoutProcess,omitempty"`
	OnFailProcess     *ProcessRequestData `json:"onFailProcess,omitempty"`
	OnCompleteProcess *ProcessRequestData `json:"onCompleteProcess,omitempty"`
	OnCancelProcess   *ProcessRequestData `json:"onCancelProcess,omitempty"`
}
