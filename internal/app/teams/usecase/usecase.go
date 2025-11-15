package usecase

import (
	"context"
	"fmt"
	"log"

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
	const op = "teams.Usecase.GetTeam"

	fail := func(code domain.ErrorCode, message string, err error) (domain.Team, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.Team{}, domain.NewError(code, message, err)
	}

	team, err := u.TeamsRepository.FetchTeamByName(ctx, teamName)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	teamMembers, err := u.UsersRepository.FetchByTeamName(ctx, teamName)
	if err != nil {
		return fail(domain.NOT_FOUND, "resource not found", err)
	}

	team.Members = teamMembers

	return team, nil
}

func (u *Usecase) AddTeam(ctx context.Context, team *domain.Team) (domain.Team, error) {
	const op = "teams.Usecase.AddTeam"

	fail := func(code domain.ErrorCode, message string, err error) (domain.Team, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.Team{}, domain.NewError(code, message, err)
	}

	err := u.TeamsRepository.CreateTeam(ctx, team)
	if err != nil {
		return fail(domain.TEAM_EXISTS, fmt.Sprintf("%s already exists", team.TeamName), err)
	}
	for _, teamMember := range team.Members {
		var user = domain.User{
			UserID:   teamMember.UserID,
			Username: teamMember.UserName,
			TeamName: team.TeamName,
			IsActive: teamMember.IsActive,
		}
		_, err = u.UsersRepository.CreateOrUpdateUser(ctx, &user)
		if err != nil {
			return fail(domain.INTERNAL, "internal server error", err)
		}
	}

	return *team, nil
}
