package dtos

import "github.com/leoscrowi/pr-assignment-service/domain"

type GetPullRequestStatsResponse struct {
	PullRequestStats []domain.PullRequestStats `json:"stats"`
}
