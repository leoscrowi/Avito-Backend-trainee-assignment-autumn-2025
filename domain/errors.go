package domain

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorCode string

const (
	TEAM_EXISTS  ErrorCode = "Team exists"
	PR_EXISTS    ErrorCode = "Pull request exists"
	PR_MERGED    ErrorCode = "Pull request merged"
	NOT_ASSIGNED ErrorCode = "Not assigned"
	NO_CANDIDATE ErrorCode = "No candidate"
	NOT_FOUND    ErrorCode = "Not found"

	INTERNAL    ErrorCode = "Internal server error"
	BAD_REQUEST ErrorCode = "Bad request"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *ErrorResponse) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

func NewError(code ErrorCode, message string, err error) *ErrorResponse {
	return &ErrorResponse{Code: code, Message: message, Err: err}
}

type APIErrorResponse struct {
	Error ErrorResponse `json:"error"`
}

func WriteError(w http.ResponseWriter, response *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeMapper(response))
	err := json.NewEncoder(w).Encode(APIErrorResponse{
		Error: ErrorResponse{
			Code:    response.Code,
			Message: response.Message,
		},
	})
	if err != nil {
		return
	}
}

func ConvertToErrorResponse(err error) *ErrorResponse {
	var errR *ErrorResponse
	if errors.As(err, &errR) {
		return &ErrorResponse{
			Code:    errR.Code,
			Message: errR.Message,
		}
	}

	return &ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: err.Error(),
	}
}
