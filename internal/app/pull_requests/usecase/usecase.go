package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests"
)

type Usecase struct {
	Repository pull_requests.Repository
}

func NewUsecase(repository pull_requests.Repository) *Usecase {
	return &Usecase{Repository: repository}
}

func (u *Usecase) ReassignPullRequest(ctx context.Context, pullRequestID uuid.UUID, newUserID uuid.UUID) (*domain.PullRequest, error) {
	// TODO: implement
	panic("implement me")
}

func (u *Usecase) MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (domain.PullRequest, error) {
	pr, err := u.Repository.MergePullRequest(ctx, pullRequestID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return pr, err
}

func (u *Usecase) CreatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error) {
	// TODO: implement
	panic("implement me")
}
