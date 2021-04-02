package alks

// BaseResponse represents basic fields included in all ALKS REST API responses
type BaseResponse struct {
	StatusMessage string   `json:"statusMessage,omitempty"`
	Errors        []string `json:"errors,omitempty"`
	RequestID     string   `json:"requestId,omitempty"`
}

// RequestFailed returns a boolean indicating if an ALKS response contained an error
func (b BaseResponse) RequestFailed() bool {
	return (b.StatusMessage != "Success" && b.StatusMessage != "") || len(b.Errors) != 0
}

// GetErrors returns a list of error messages from an ALKS response
func (b BaseResponse) GetErrors() []string {
	var errorMessages []string
	errorMessages = append(errorMessages, b.Errors...)

	if len(errorMessages) == 0 {
		errorMessages = []string{
			b.StatusMessage,
		}
	}

	return errorMessages
}
