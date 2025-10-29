package common

type StepConfiguration struct {
	ErrorPolicyName          *string `json:"errorPolicyName,omitempty"`
	UseCache                 bool    `json:"useCache"`
	CacheKey                 *string `json:"cacheKey,omitempty"`
	CacheTimeoutSeconds      int     `json:"cacheTimeoutSeconds"`
	Stream                   bool    `json:"stream"`
	StreamAllStatuses        bool    `json:"streamAllStatuses"`
	PreExecutionScript       *string `json:"preExecutionScript,omitempty"`
	ConstructInputFromScript bool    `json:"constructInputFromScript"`
	ScriptLanguage           *string `json:"scriptLanguage,omitempty"`
	JQTransformer            *string `json:"jqTransformer,omitempty"`
	RateLimitMaxRequests     int     `json:"rateLimitMaxRequests"`
	RateLimitWindowSeconds   int     `json:"rateLimitWindowSeconds"`
}
