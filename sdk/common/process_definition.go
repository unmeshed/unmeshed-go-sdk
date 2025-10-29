package common

type ProcessDefinition struct {
	OrgID         *int                         `json:"orgId,omitempty"`
	Namespace     *string                      `json:"namespace,omitempty"`
	Name          *string                      `json:"name,omitempty"`
	Version       *int                         `json:"version,omitempty"`
	Type          *ProcessType                 `json:"type,omitempty"`
	Description   *string                      `json:"description,omitempty"`
	CreatedBy     *string                      `json:"createdBy,omitempty"`
	UpdatedBy     *string                      `json:"updatedBy,omitempty"`
	Created       int                          `json:"created"`
	Updated       int                          `json:"updated"`
	Configuration *ProcessConfiguration        `json:"configuration,omitempty"`
	Steps         []*StepDefinition            `json:"steps,omitempty"`
	DefaultInput  map[string]interface{}       `json:"defaultInput,omitempty"`
	DefaultOutput map[string]interface{}       `json:"defaultOutput,omitempty"`
	OutputMapping map[string]interface{}       `json:"outputMapping,omitempty"`
	Metadata      map[string][]*StepDependency `json:"metadata,omitempty"`
	Tags          []*TagValue                  `json:"tags,omitempty"`
}
