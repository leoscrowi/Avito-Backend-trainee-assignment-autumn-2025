package users

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type UseCase interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) error
	GetReview(ctx context.Context, userID string) ([]domain.PullRequestShort, error)
}
