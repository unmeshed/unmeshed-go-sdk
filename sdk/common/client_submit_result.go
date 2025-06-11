package common

type ClientSubmitResult struct {
	ProcessID    int64
	StepID       int64
	ErrorMessage string
	StatusCode   int
}

func NewClientSubmitResult(processID int64, stepID int64, statusCode int, errorMessage string) *ClientSubmitResult {
	return &ClientSubmitResult{
		ProcessID:    processID,
		StepID:       stepID,
		ErrorMessage: errorMessage,
		StatusCode:   statusCode,
	}
}

func (c *ClientSubmitResult) GetErrorMessage() string {
	return c.ErrorMessage
}
