package models

// swagger:model
type APIResponse struct {
	// The response message
	// example: Success
	Message string `json:"message"`

	// Optional error
	// example: Something went wrong
	Error string `json:"error,omitempty"`

	// Optional data payload
	Data interface{} `json:"data,omitempty"`
}
