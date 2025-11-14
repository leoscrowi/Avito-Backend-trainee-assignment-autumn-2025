package pull_requests

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Repository interface {
	CreatePullRequest(ctx context.Context, pr *domain.PullRequest) error
	MergePullRequest(ctx context.Context, prID string) (domain.PullRequest, error)

	GetReviewersID(ctx context.Context, prID string) ([]string, error)
	DeleteReviewer(ctx context.Context, prID, reviewerID string) error
	AddReviewer(ctx context.Context, prID, reviewerID string) error

	FetchByID(ctx context.Context, prID string) (domain.PullRequest, error)
	IsMerged(ctx context.Context, prID string) (bool, error)

	FindPullRequestsByUserID(ctx context.Context, userID string) ([]domain.PullRequestShort, error)
}
