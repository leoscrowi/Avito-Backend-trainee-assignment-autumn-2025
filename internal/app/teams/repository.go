package teams

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Repository interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	FetchByName(ctx context.Context, teamName string) (domain.Team, error)
}
