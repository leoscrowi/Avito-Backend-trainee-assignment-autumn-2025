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

func (u *Usecase) SetIsActive(ctx context.Context, userID uuid.UUID, isActive bool) error {
	// TODO: implement
	panic("implement me")
}

func (u *Usecase) GetReview(userID uuid.UUID) []domain.PullRequest {
	// TODO: implement
	panic("implement me")
}
