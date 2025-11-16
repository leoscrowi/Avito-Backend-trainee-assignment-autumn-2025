package domain

import (
	"time"
)

type Status string

const (
	OPEN   Status = "OPEN"
	MERGED Status = "MERGED"
)

type PullRequest struct {
	PullRequestID     string    `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName   string    `json:"pull_request_name" db:"pull_request_name"`
	AuthorID          string    `json:"author_id" db:"author_id"`
	Status            Status    `json:"status" db:"status"`
	AssignedReviewers []string  `json:"assigned_reviewers"`
	NeedMoreReviewers bool      `json:"need_more_reviewers" db:"need_more_reviewers"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	MergedAt          time.Time `json:"merged_at" db:"merged_at"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorID        string `json:"author_id" db:"author_id"`
	Status          Status `json:"status" db:"status"`
}
