package usecase

import (
	"context"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/users"
)

type Usecase struct {
	Repository users.Repository
}

func NewUsecase(repository users.Repository) *Usecase {
	return &Usecase{Repository: repository}
}

func (u *Usecase) GetTeam(ctx context.Context, teamName string) (domain.Team, error) {
	// TODO: implement
	panic("implement me")
}

func (u *Usecase) AddTeam(ctx context.Context, team *domain.Team) (domain.Team, error) {
	// TODO: implement
	panic("implement me")
}
