package usecase

import (
	"context"
	"github.com/leoscrowi/pr-assignment-service/internal/app/teams"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/users"
)

type Usecase struct {
	UsersRepository users.Repository
	TeamsRepository teams.Repository
}

func NewUsecase(uRepository users.Repository, tRepository teams.Repository) *Usecase {
	return &Usecase{UsersRepository: uRepository, TeamsRepository: tRepository}
}

func (u *Usecase) GetTeam(ctx context.Context, teamName string) (domain.Team, error) {
	team, err := u.TeamsRepository.FetchByName(ctx, teamName)
	if err != nil {
		return domain.Team{}, err
	}

	teamMembers, err := u.UsersRepository.FetchByTeamName(ctx, teamName)
	if err != nil {
		return domain.Team{}, err
	}
	team.Members = teamMembers

	return team, nil
}

func (u *Usecase) AddTeam(ctx context.Context, team *domain.Team) (domain.Team, error) {
	err := u.TeamsRepository.CreateTeam(ctx, team)
	if err != nil {
		return domain.Team{}, err
	}
	for _, teamMember := range team.Members {
		var user = domain.User{
			UserID:   teamMember.UserID,
			Username: teamMember.UserName,
			TeamName: team.TeamName,
			IsActive: true,
		}
		_, err = u.UsersRepository.CreateOrUpdateUser(ctx, &user)
		if err != nil {
			return domain.Team{}, err
		}
	}

	return *team, nil
}
