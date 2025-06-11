package common

type WorkResponseTracker struct {
	WorkResponse  *WorkResponse
	RetryCount    int
	QueuedTime    int64
	StepPollState *StepPollState
}

func NewWorkResponseTracker(workResponse *WorkResponse) *WorkResponseTracker {
	return &WorkResponseTracker{
		WorkResponse:  workResponse,
		RetryCount:    0,
		QueuedTime:    0,
		StepPollState: nil,
	}
}
