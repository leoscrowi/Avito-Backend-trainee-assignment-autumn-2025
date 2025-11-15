package dtos

import "github.com/leoscrowi/pr-assignment-service/domain"

type GetTeamRequest struct {
	TeamName string `json:"team_name"`
}

type AddTeamResponse struct {
	Team domain.Team `json:"team"`
}
