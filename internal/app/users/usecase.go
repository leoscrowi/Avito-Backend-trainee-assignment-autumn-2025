package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoscrowi/pr-assignment-service/domain"
)

type UseCase interface {
	SetIsActive(ctx context.Context, userID uuid.UUID, isActive bool) error
	GetReview(userID uuid.UUID) []domain.PullRequest
}
