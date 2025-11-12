package domain

import (
	"github.com/google/uuid"
	"time"
)

type Status int

const (
	OPEN Status = iota
	MERGED
)

type PullRequest struct {
	PullRequestId     uuid.UUID
	PullRequestName   string
	AuthorId          uuid.UUID
	Status            Status
	AssignedReviewers []uuid.UUID
	NeedMoreReviewers bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type PullRequestShort struct {
	PullRequestId   uuid.UUID
	PullRequestName string
	AuthorId        uuid.UUID
	Status          Status
}
