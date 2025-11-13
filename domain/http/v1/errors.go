package v1

type ErrorCode string

const (
	TEAM_EXISTS  ErrorCode = "Team exists"
	PR_EXISTS    ErrorCode = "Pull request exists"
	PR_MERGED    ErrorCode = "Pull request merged"
	NOT_ASSIGNED ErrorCode = "Not assigned"
	NO_CANDIDATE ErrorCode = "No candidate"
	NOT_FOUND    ErrorCode = "Not found"
)

type ErrorResponse struct {
	Error ErrorInfo
}

type ErrorInfo struct {
	Code    ErrorCode
	Message string
}
