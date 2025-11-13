package pull_requests

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Repository interface {
	CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (domain.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pr *domain.PullRequest) (domain.PullRequest, error)
}
