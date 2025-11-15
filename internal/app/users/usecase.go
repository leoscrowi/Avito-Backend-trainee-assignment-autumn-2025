package users

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Usecase interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error)
	GetReview(ctx context.Context, userID string) ([]domain.PullRequestShort, error)
}
