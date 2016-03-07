package kii

import (
	"fmt"
	"encoding/json"
)

// ErrorResponse represents error response returned by Kii Cloud.
type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message string `json:"message"`
}

// CloudError represents error returned by Kii Cloud.
type CloudError struct {
	ErrorResponse
	HTTPStatus int
	RawResponse string
}

func newCloudError(httpStatus int, rawResponse []byte) *CloudError {
	var ce CloudError
	ce.HTTPStatus = httpStatus
	ce.RawResponse = string(rawResponse)
	json.Unmarshal(rawResponse, &ce)
	return &ce
}

func (e CloudError) Error() string {
	return fmt.Sprintf("%s : %s (%d)", e.ErrorCode, e.Message, e.HTTPStatus)
}
