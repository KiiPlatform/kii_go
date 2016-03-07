package kii

import (
	"fmt"
	"encoding/json"
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
	RawResponse string
}

func newCloudError(httpStatus int, rawResponse []byte) *CloudError {
	var ce CloudError
	ce.HttpStatus = httpStatus
	ce.RawResponse = string(rawResponse)
	json.Unmarshal(rawResponse, &ce)
	return &ce
}

func (e CloudError) Error() string {
	return fmt.Sprintf("%s : %s (%d)", e.ErrorCode, e.Message, e.HttpStatus)
}
