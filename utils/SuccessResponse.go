package utils

// SuccessResponse represents a successful response from the API
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
