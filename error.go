package kii

import (
	"fmt"
)

// Represents Error Response Returned by Kii Cloud.
type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message string `json:"message"`
}

// Represents Error returned by Kii Cloud.
type CloudError struct {
	ErrorResponse
	HttpStatus int
}

func (e CloudError) Error() string {
	return fmt.Sprintf("%s : %s (%d)", e.ErrorCode, e.Message, e.HttpStatus)
}
