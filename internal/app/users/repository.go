package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Repository interface {
	SetIsActive(ctx context.Context, userID uuid.UUID, isActive bool) error
	CreateOrUpdateUser(ctx context.Context, user *domain.User) (uuid.UUID, error)
}
