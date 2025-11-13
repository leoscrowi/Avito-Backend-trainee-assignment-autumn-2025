package pull_requests

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Repository interface {
	MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (domain.PullRequest, error)
	CreatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error)
	UpdatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error)
}