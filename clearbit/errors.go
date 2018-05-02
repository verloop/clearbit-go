package clearbit

import (
	"encoding/json"
	"fmt"
)

// ApiError represents a Clearbit API Error response
// https://clearbit.com/docs#errors
type ApiError struct {
	Errors []ErrorDetail `json:"error"`
}
type apiError = ApiError

// ErrorDetail represents an individual item in an apiError.
type ErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ErrorDetailWrapper is used for single error
type ErrorDetailWrapper struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents an individual item in an apiError.
func (e ApiError) Error() string {
	if len(e.Errors) > 0 {
		err := e.Errors[0]
		return fmt.Sprintf("clearbit: %s %v", err.Type, err.Message)
	}
	return ""
}

// UnmarshalJSON is used to be able to read dynamic json
//
// This is because sometimes our errors are not arrays of ErrorDetail but a
// single ErrorDetail
func (e *ApiError) UnmarshalJSON(b []byte) (err error) {
	errorDetail, errors := ErrorDetailWrapper{}, []ErrorDetail{}
	if err = json.Unmarshal(b, &errors); err == nil {
		e.Errors = errors
		return
	}

	if err = json.Unmarshal(b, &errorDetail); err == nil {
		errors = append(errors, errorDetail.Error)
		e.Errors = errors
		return
	}

	return err
}

// Empty returns true if empty. Otherwise, at least 1 error message/code is
// present and false is returned.
func (e *ApiError) Empty() bool {
	if len(e.Errors) == 0 {
		return true
	}
	return false
}

// relevantError returns any non-nil http-related error (creating the request,
// getting the response, decoding) if any. If the decoded apiError is non-zero
// the apiError is returned. Otherwise, no errors occurred, returns nil.
func relevantError(httpError error, ae apiError) error {
	if httpError != nil {
		return httpError
	}
	if ae.Empty() {
		return nil
	}
	return ae
}
