package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/teams"
	"github.com/leoscrowi/pr-assignment-service/internal/app/teams/dtos"
	"github.com/leoscrowi/pr-assignment-service/internal/utils"
)

type TeamsController struct {
	usecase teams.Usecase
}

func NewTeamsController(usecase teams.Usecase) *TeamsController {
	return &TeamsController{usecase: usecase}
}

func (c *TeamsController) GetTeam(w http.ResponseWriter, r *http.Request) {
	var req dtos.GetTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", err))
		return
	}

	if req.TeamName == "" {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", fmt.Errorf("wrong json format")))
		return
	}

	team, err := c.usecase.GetTeam(r.Context(), req.TeamName)
	if err != nil {
		domain.WriteError(w, domain.ConvertToErrorResponse(err))
		return
	}
	utils.WriteHeader(w, http.StatusOK, &team)
}

func (c *TeamsController) AddTeam(w http.ResponseWriter, r *http.Request) {
	var team domain.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", err))
		return
	}

	if team.TeamName == "" || len(team.Members) == 0 {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", fmt.Errorf("wrong json format")))
		return
	}

	newTeam, err := c.usecase.AddTeam(r.Context(), &team)
	if err != nil {
		domain.WriteError(w, domain.ConvertToErrorResponse(err))
	}

	var resp = dtos.AddTeamResponse{Team: newTeam}
	utils.WriteHeader(w, http.StatusCreated, &resp)
}
