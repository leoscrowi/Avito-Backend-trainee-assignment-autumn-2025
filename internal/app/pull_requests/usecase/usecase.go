package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/users"
)

type Usecase struct {
	Repository users.Repository
}

func NewUsecase(repository users.Repository) *Usecase {
	return &Usecase{Repository: repository}
}

func (u *Usecase) ReassignPullRequest(ctx context.Context, pullRequestID uuid.UUID, newUserID uuid.UUID) (*domain.PullRequest, error) {
	// TODO: implement
	panic("implement me")
}

func (u *Usecase) MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (domain.PullRequest, error) {
	// TODO: implement
	panic("implement me")
}

func (u *Usecase) CreatePullRequest(ctx context.Context, pullRequest *domain.PullRequest) (domain.PullRequest, error) {
	// TODO: implement
	panic("implement me")
}
