package domain

type ErrorCode string

const (
	TEAM_EXISTS  ErrorCode = "Team exists"
	PR_EXISTS    ErrorCode = "Pull request exists"
	PR_MERGED    ErrorCode = "Pull request merged"
	NOT_ASSIGNED ErrorCode = "Not assigned"
	NO_CANDIDATE ErrorCode = "No candidate"
	NOT_FOUND    ErrorCode = "Not found"

	INTERNAL ErrorCode = "Internal server error"
)

type ErrorResponse struct {
	Code    ErrorCode
	Message string
	Err     error
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
