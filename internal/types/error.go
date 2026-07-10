package types

import "fmt"

type APIError struct {
	Code    int               `json:"statusCode"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%d]: %s", e.Code, e.Message)
}

func NewAPIError(code int, err error) APIError {
	return APIError{
		Code:    code,
		Message: err.Error(),
	}
}

func NewAPIErrorWithDetail(code int, err error, errors map[string]string) APIError {
	return APIError{
		Code:    code,
		Message: err.Error(),
		Errors:  errors,
	}
}
