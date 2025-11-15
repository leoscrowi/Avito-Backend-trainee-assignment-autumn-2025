package dtos

import "github.com/leoscrowi/pr-assignment-service/domain"

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	User domain.User `json:"user"`
}

type GetReviewRequest struct {
	UserID string `json:"user_id"`
}

type GetReviewResponse struct {
	UserID       string                    `json:"user_id"`
	PullRequests []domain.PullRequestShort `json:"pull_requests"`
}
