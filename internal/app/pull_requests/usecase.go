package pull_requests

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Usecase interface {
	ReassignPullRequest(ctx context.Context, pullRequestID string, oldUserID string) (domain.PullRequest, string, error)
	MergePullRequest(ctx context.Context, pullRequestID string) (domain.PullRequest, error)
	CreatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error)
}
