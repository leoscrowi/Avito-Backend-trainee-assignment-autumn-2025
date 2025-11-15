package usecase

import (
	"context"
	"log"

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

func (u *Usecase) SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	const op = "users.Usecase.SetIsActive"

	fail := func(code domain.ErrorCode, message string, err error) (domain.User, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.User{}, domain.NewError(code, message, err)
	}

	user, err := u.UsersRepository.FetchByID(ctx, userID)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	err = u.UsersRepository.SetIsActive(ctx, userID, isActive)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	user.IsActive = true
	return user, nil
}

func (u *Usecase) GetReview(ctx context.Context, userID string) ([]domain.PullRequestShort, error) {
	const op = "users.Usecase.GetReview"

	fail := func(code domain.ErrorCode, message string, err error) ([]domain.PullRequestShort, error) {
		log.Printf("%s: %v\n", op, err)
		return []domain.PullRequestShort{}, domain.NewError(code, message, err)
	}

	prs, err := u.PullRequestsRepository.FindPullRequestsByUserID(ctx, userID)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	return prs, nil
}
