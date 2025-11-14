package pull_requests

import (
	"context"

	"github.com/google/uuid"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Repository interface {
	CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (domain.PullRequest, error)
	FindPullRequestsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.PullRequestShort, error)

	// TODO: ??? подумать как можно сделать
	ReassignPullRequest(ctx context.Context, prID uuid.UUID, newReviewerID uuid.UUID) (uuid.UUID, error)
	MergePullRequest(ctx context.Context, prID uuid.UUID) (domain.PullRequest, error)
}
