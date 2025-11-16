package domain

func statusCodeMapper(err *ErrorResponse) int {
	switch err.Code {
	case TEAM_EXISTS:
		return 400
	case NOT_FOUND:
		return 404
	case PR_MERGED:
		return 409
	case NOT_ASSIGNED:
		return 409
	case NO_CANDIDATE:
		return 409
	case BAD_REQUEST:
		return 400
	case UNAUTHORIZED:
		return 401
	case PR_EXISTS:
		return 409
	default:
		return 500
	}
}
