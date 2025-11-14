package users

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Repository interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) error
	CreateOrUpdateUser(ctx context.Context, user *domain.User) (string, error)
	FetchByID(ctx context.Context, userID string) (domain.User, error)
	FetchByTeamName(ctx context.Context, teamName string) ([]domain.TeamMember, error)

	GetActiveUsersIDByTeam(ctx context.Context, teamName string) ([]string, error)
}
