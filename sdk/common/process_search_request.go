package common

import (
	"time"
)

type ProcessSearchRequest struct {
	StartTimeEpoch int64                `json:"startTimeEpoch"`
	EndTimeEpoch   *int64               `json:"endTimeEpoch,omitempty"`
	Namespace      *string              `json:"namespace,omitempty"`
	ProcessTypes   []ProcessType        `json:"processTypes,omitempty"`
	TriggerTypes   []ProcessTriggerType `json:"triggerTypes,omitempty"`
	Names          []string             `json:"names,omitempty"`
	ProcessIds     []int64              `json:"processIds,omitempty"`
	CorrelationIds []string             `json:"correlationIds,omitempty"`
	RequestIds     []string             `json:"requestIds,omitempty"`
	Statuses       []ProcessStatus      `json:"statuses,omitempty"`
	Limit          int                  `json:"limit"`
	Offset         int                  `json:"offset"`
}

func NewProcessSearchRequest() *ProcessSearchRequest {
	return &ProcessSearchRequest{
		StartTimeEpoch: int64(time.Now().UnixMilli()) - (60 * 1000 * 60 * 24),
		Limit:          10,
		Offset:         0,
	}
}
