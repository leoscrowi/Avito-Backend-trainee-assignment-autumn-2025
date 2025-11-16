package v1

import (
	"net/http"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/stats"
	"github.com/leoscrowi/pr-assignment-service/internal/app/stats/dtos"
	"github.com/leoscrowi/pr-assignment-service/internal/utils"
)

type StatsController struct {
	usecase stats.Usecase
}

func NewStatsController(usecase stats.Usecase) *StatsController {
	return &StatsController{usecase: usecase}
}

func (c *StatsController) GetPullRequestStats(w http.ResponseWriter, r *http.Request) {
	st, err := c.usecase.GetPullRequestStats(r.Context())
	if err != nil {
		domain.WriteError(w, domain.ConvertToErrorResponse(err))
		return
	}

	var resp = dtos.GetPullRequestStatsResponse{
		PullRequestStats: st,
	}
	utils.WriteHeader(w, http.StatusOK, &resp)
}
