package stats

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Usecase interface {
	GetPullRequestStats(ctx context.Context) ([]domain.PullRequestStats, error)
}
