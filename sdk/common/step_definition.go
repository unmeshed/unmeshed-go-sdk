package common

type StepDefinition struct {
	OrgID         *int                   `json:"orgId,omitempty"`
	Namespace     *string                `json:"namespace,omitempty"`
	Name          *string                `json:"name,omitempty"`
	Type          *StepType              `json:"type,omitempty"`
	Ref           *string                `json:"ref,omitempty"`
	Optional      bool                   `json:"optional"`
	CreatedBy     *string                `json:"createdBy,omitempty"`
	UpdatedBy     *string                `json:"updatedBy,omitempty"`
	Description   *string                `json:"description,omitempty"`
	Label         *string                `json:"label,omitempty"`
	Created       int                    `json:"created"`
	Updated       int                    `json:"updated"`
	Configuration *StepConfiguration     `json:"configuration,omitempty"`
	Children      []*StepDefinition      `json:"children,omitempty"`
	Input         map[string]interface{} `json:"input,omitempty"`
	Output        map[string]interface{} `json:"output,omitempty"`
}
