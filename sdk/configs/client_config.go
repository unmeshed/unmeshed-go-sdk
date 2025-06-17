package configs

import "github.com/unmeshed/unmeshed-go-sdk/sdk/common"

type ClientConfig struct {
	Namespace                      string
	BaseURL                        string
	Port                           int
	ConnectionTimeoutSecs          int64
	SubmitClientPollTimeoutSeconds float64
	StepTimeoutMillis              int64
	DelayMillis                    int64
	WorkRequestBatchSize           int
	StepSubmissionAttempts         int64
	ClientID                       string
	AuthToken                      string
	MaxWorkers                     int64
	PollRequestData                common.PollRequestData
	ResponseSubmitBatchSize        int
	permanentErrorKeywords         []string
	MaxSubmitAttempts              int64
}

func NewClientConfig() *ClientConfig {
	defaultNamespace := "default"
	defaultBaseURL := "http://localhost"
	defaultPort := 8080
	defaultConnectionTimeoutSecs := int64(60)
	defaultSubmitClientPollTimeoutSeconds := float64(30)
	defaultStepTimeoutMillis := int64(5000)
	defaultDelayMillis := int64(100)
	defaultWorkRequestBatchSize := int(100)
	defaultStepSubmissionAttempts := int64(3)
	defaultMaxWorkers := int64(20)
	responseSubmitBatchSize := int(500)

	return &ClientConfig{
		Namespace:                      defaultNamespace,
		BaseURL:                        defaultBaseURL,
		Port:                           defaultPort,
		ConnectionTimeoutSecs:          defaultConnectionTimeoutSecs,
		SubmitClientPollTimeoutSeconds: defaultSubmitClientPollTimeoutSeconds,
		StepTimeoutMillis:              defaultStepTimeoutMillis,
		DelayMillis:                    defaultDelayMillis,
		WorkRequestBatchSize:           defaultWorkRequestBatchSize,
		StepSubmissionAttempts:         defaultStepSubmissionAttempts,
		MaxWorkers:                     defaultMaxWorkers,
		ResponseSubmitBatchSize:        responseSubmitBatchSize,
		permanentErrorKeywords: []string{
			"Invalid request, step is not in RUNNING state",
			"please poll the latest and update",
		},
	}
}

func (c *ClientConfig) PermanentErrorKeywords() []string {
	return c.permanentErrorKeywords
}

func (c *ClientConfig) HasToken() bool {
	return c.AuthToken != "" // No need for nil check as it's a string now
}

func (c *ClientConfig) GetNamespace() string            { return c.Namespace }
func (c *ClientConfig) GetBaseURL() string              { return c.BaseURL }
func (c *ClientConfig) GetPort() int                    { return c.Port }
func (c *ClientConfig) GetConnectionTimeoutSecs() int64 { return c.ConnectionTimeoutSecs }
func (c *ClientConfig) GetSubmitClientPollTimeoutSeconds() float64 {
	return c.SubmitClientPollTimeoutSeconds
}
func (c *ClientConfig) GetStepTimeoutMillis() int64                { return c.StepTimeoutMillis }
func (c *ClientConfig) GetDelayMillis() int64                      { return c.DelayMillis }
func (c *ClientConfig) GetWorkRequestBatchSize() int               { return c.WorkRequestBatchSize }
func (c *ClientConfig) GetStepSubmissionAttempts() int64           { return c.StepSubmissionAttempts }
func (c *ClientConfig) GetClientID() string                        { return c.ClientID }
func (c *ClientConfig) GetAuthToken() string                       { return c.AuthToken }
func (c *ClientConfig) GetMaxWorkers() int64                       { return c.MaxWorkers }
func (c *ClientConfig) GetPollRequestData() common.PollRequestData { return c.PollRequestData }
func (c *ClientConfig) GetResponseSubmitBatchSize() int            { return c.ResponseSubmitBatchSize }
func (c *ClientConfig) GetMaxSubmitAttempts() int64                { return c.MaxSubmitAttempts }

func (c *ClientConfig) SetNamespace(namespace string) {
	if namespace == "" {
		panic("namespace cannot be empty")
	}
	c.Namespace = namespace
}

func (c *ClientConfig) SetBaseURL(baseURL string) {
	if baseURL == "" {
		panic("Base URL cannot be empty")
	}
	c.BaseURL = baseURL
}

func (c *ClientConfig) SetPort(port int) {
	if port <= 0 {
		panic("Port number must be a positive integer")
	}
	c.Port = port
}

func (c *ClientConfig) SetConnectionTimeoutSecs(connectionTimeoutSecs int64) {
	if connectionTimeoutSecs <= 0 {
		panic("Connection timeout must be a positive integer")
	}
	c.ConnectionTimeoutSecs = connectionTimeoutSecs
}

func (c *ClientConfig) SetSubmitClientPollTimeoutSeconds(submitClientPollTimeoutSecs float64) {
	if submitClientPollTimeoutSecs <= 0 {
		panic("Submit client poll timeouut must be a positive integer")
	}
	c.SubmitClientPollTimeoutSeconds = submitClientPollTimeoutSecs
}

func (c *ClientConfig) SetStepTimeoutMillis(stepTimeoutMillis int64) {
	if stepTimeoutMillis <= 0 {
		panic("Step timeout must be a positive integer")
	}
	c.StepTimeoutMillis = stepTimeoutMillis
}

func (c *ClientConfig) SetDelayMillis(delayMillis int64) {
	if delayMillis < 0 {
		panic("Delay cannot be negative")
	}
	c.DelayMillis = delayMillis
}

func (c *ClientConfig) SetWorkRequestBatchSize(workRequestBatchSize int) {
	if workRequestBatchSize <= 0 {
		panic("Work request batch size must be a positive integer")
	}
	c.WorkRequestBatchSize = workRequestBatchSize
}

func (c *ClientConfig) SetStepSubmissionAttempts(stepSubmissionAttempts int64) {
	if stepSubmissionAttempts <= 0 {
		panic("Step submission attempts must be a positive integer")
	}
	c.StepSubmissionAttempts = stepSubmissionAttempts
}

func (c *ClientConfig) SetClientID(clientID string) {
	c.ClientID = clientID
}

func (c *ClientConfig) SetAuthToken(authToken string) {
	c.AuthToken = authToken
}

func (c *ClientConfig) SetMaxWorkers(maxWorkers int64) {
	if maxWorkers <= 0 {
		panic("Max Workers count must be a positive integer")
	}
	c.MaxWorkers = maxWorkers
}

func (c *ClientConfig) SetPollRequestData(pollRequestData common.PollRequestData) {
	c.PollRequestData = pollRequestData
}

func (c *ClientConfig) SetResponseSubmitBatchSize(responseSubmitBatchSize int) {
	if responseSubmitBatchSize <= 0 {
		panic("Response Submit Batch Size must be a positive integer")
	}
	c.ResponseSubmitBatchSize = responseSubmitBatchSize
}

func (c *ClientConfig) SetMaxSubmitAttempts(maxSubmitAttempts int64) {
	if maxSubmitAttempts <= 0 {
		panic("Max submit attempts must be a positive integer")
	}
	c.MaxSubmitAttempts = maxSubmitAttempts
}
