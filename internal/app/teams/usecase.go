package teams

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Usecase interface {
	GetTeam(ctx context.Context, teamName string) (domain.Team, error)
	AddTeam(ctx context.Context, team *domain.Team) (domain.Team, error)
}
