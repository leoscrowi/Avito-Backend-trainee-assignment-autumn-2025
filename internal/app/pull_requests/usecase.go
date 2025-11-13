package pull_requests

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Usecase interface {
	ReassignPullRequest(ctx context.Context, pullRequestID uuid.UUID, newUserID uuid.UUID) (*domain.PullRequest, error)
	MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (domain.PullRequest, error)
	CreatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error)
}
