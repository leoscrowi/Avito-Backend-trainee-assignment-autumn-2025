package users

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	SetIsActive(ctx context.Context, userID uuid.UUID, isActive bool) error
}
