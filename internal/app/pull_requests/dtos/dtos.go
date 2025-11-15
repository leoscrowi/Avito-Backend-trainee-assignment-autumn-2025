package dtos

import "github.com/leoscrowi/pr-assignment-service/domain"

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePRResponse struct {
	PR domain.PullRequest `json:"pr"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type MergePRResponse struct {
	PR domain.PullRequest `json:"pr"`
}

type ReassignPRRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type ReassignPRResponse struct {
	PR         domain.PullRequest `json:"pr"`
	ReplacedBy string             `json:"replaced_by"`
}
