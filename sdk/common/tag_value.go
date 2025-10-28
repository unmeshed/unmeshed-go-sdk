package common

// TagValue represents a name-value pair for tagging entities.
type TagValue struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Of creates a new TagValue, similar to TagValue.of() in Java.
func Of(name, value string) *TagValue {
	return &TagValue{
		Name:  name,
		Value: value,
	}
}

// GetName returns the name of the TagValue.
func (t *TagValue) GetName() string {
	return t.Name
}

// GetValue returns the value of the TagValue.
func (t *TagValue) GetValue() string {
	return t.Value
}
