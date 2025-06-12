package common

type ProcessActionResponseDetailData struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewProcessActionResponseDetailData(id, message, error string) *ProcessActionResponseDetailData {
	return &ProcessActionResponseDetailData{
		ID:      id,
		Message: message,
		Error:   error,
	}
}

func (p *ProcessActionResponseDetailData) GetID() string {
	return p.ID
}

func (p *ProcessActionResponseDetailData) GetMessage() string {
	return p.Message
}

func (p *ProcessActionResponseDetailData) GetError() string {
	return p.Error
}

func (p *ProcessActionResponseDetailData) SetID(id string) {
	p.ID = id
}

func (p *ProcessActionResponseDetailData) SetMessage(message string) {
	p.Message = message
}

func (p *ProcessActionResponseDetailData) SetError(error string) {
	p.Error = error
}

type ProcessActionResponseData struct {
	Count   int                                `json:"count"`
	Details []*ProcessActionResponseDetailData `json:"details,omitempty"`
}

func NewProcessActionResponseData() *ProcessActionResponseData {
	return &ProcessActionResponseData{
		Count:   0,
		Details: make([]*ProcessActionResponseDetailData, 0),
	}
}

func (p *ProcessActionResponseData) GetCount() int {
	return p.Count
}

func (p *ProcessActionResponseData) GetDetails() []*ProcessActionResponseDetailData {
	return p.Details
}

func (p *ProcessActionResponseData) SetCount(count int) {
	p.Count = count
}

func (p *ProcessActionResponseData) SetDetails(details []*ProcessActionResponseDetailData) {
	p.Details = details
}

// FromMap creates a ProcessActionResponseData from a map
func FromMap(data map[string]interface{}) (*ProcessActionResponseData, error) {
	result := NewProcessActionResponseData()

	// Set count
	if count, ok := data["count"].(float64); ok {
		result.Count = int(count)
	}

	// Set details
	if details, ok := data["details"].([]interface{}); ok {
		for _, item := range details {
			if detailMap, ok := item.(map[string]interface{}); ok {
				detail := &ProcessActionResponseDetailData{}

				if id, ok := detailMap["id"].(string); ok {
					detail.ID = id
				}
				if message, ok := detailMap["message"].(string); ok {
					detail.Message = message
				}
				if error, ok := detailMap["error"].(string); ok {
					detail.Error = error
				}

				result.Details = append(result.Details, detail)
			}
		}
	}

	return result, nil
}
