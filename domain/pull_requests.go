package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	OPEN   Status = "OPEN"
	MERGED Status = "MERGED"
)

// TODO: кажется, можно добавить nullable прям в кавычках
type PullRequest struct {
	PullRequestID     uuid.UUID   `json:"pull_request_id"`
	PullRequestName   string      `json:"pull_request_name"`
	AuthorId          uuid.UUID   `json:"author_id"`
	Status            Status      `json:"status"`
	AssignedReviewers []uuid.UUID `json:"assigned_reviewers"`
	NeedMoreReviewers bool        `json:"need_more_reviewers"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type PullRequestShort struct {
	PullRequestID   uuid.UUID `json:"pull_request_id"`
	PullRequestName string    `json:"pull_request_name"`
	AuthorID        uuid.UUID `json:"author_id"`
	Status          Status    `json:"status"`
}
