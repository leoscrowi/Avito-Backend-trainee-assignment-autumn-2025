package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests"
	"github.com/leoscrowi/pr-assignment-service/internal/app/users"
)

type Usecase struct {
	UsersRepository        users.Repository
	PullRequestsRepository pull_requests.Repository
}

func NewUsecase(uRepository users.Repository, prRepository pull_requests.Repository) *Usecase {
	return &Usecase{UsersRepository: uRepository, PullRequestsRepository: prRepository}
}

func (u *Usecase) SetIsActive(ctx context.Context, userID uuid.UUID, isActive bool) (domain.User, error) {
	const op = "users.Usecase.SetIsActive"

	user, err := u.UsersRepository.FetchByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: user not found: %v", op, err)
	}

	err = u.UsersRepository.SetIsActive(ctx, userID, isActive)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %v", op, err)
	}

	user.IsActive = true
	return user, nil
}

func (u *Usecase) GetReview(ctx context.Context, userID uuid.UUID) ([]domain.PullRequestShort, error) {
	const op = "users.Usecase.GetReview"

	prs, err := u.PullRequestsRepository.FindPullRequestsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return prs, nil
}
