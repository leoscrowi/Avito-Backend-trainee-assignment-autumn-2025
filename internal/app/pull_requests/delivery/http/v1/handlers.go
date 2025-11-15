package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests"
	"github.com/leoscrowi/pr-assignment-service/internal/app/pull_requests/dtos"
	"github.com/leoscrowi/pr-assignment-service/internal/utils"
)

type PullRequestController struct {
	usecase pull_requests.Usecase
}

func NewPullRequestController(usecase pull_requests.Usecase) *PullRequestController {
	return &PullRequestController{usecase: usecase}
}

func (c *PullRequestController) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", err))
		return
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", fmt.Errorf("wrong json format")))
		return
	}

	pullRequest := &domain.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}

	pr, err := c.usecase.CreatePullRequest(r.Context(), pullRequest)
	if err != nil {
		domain.WriteError(w, domain.ConvertToErrorResponse(err))
		return
	}

	var resp = dtos.CreatePRResponse{
		PR: pr,
	}
	utils.WriteHeader(w, http.StatusCreated, &resp)
}

func (c *PullRequestController) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {
	var req dtos.ReassignPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", err))
		return
	}

	if req.OldUserID == "" || req.PullRequestID == "" {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", nil))
		return
	}

	pr, replacedBy, err := c.usecase.ReassignPullRequest(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		domain.WriteError(w, domain.ConvertToErrorResponse(err))
		return
	}

	var resp = dtos.ReassignPRResponse{
		PR:         pr,
		ReplacedBy: replacedBy,
	}
	utils.WriteHeader(w, http.StatusOK, &resp)
}

func (c *PullRequestController) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req dtos.ReassignPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", err))
		return
	}

	if req.PullRequestID == "" {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "pull request ID is required", nil))
		return
	}

	pr, err := c.usecase.MergePullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		domain.WriteError(w, domain.ConvertToErrorResponse(err))
		return
	}

	var resp = dtos.MergePRResponse{
		PR: pr,
	}
	utils.WriteHeader(w, http.StatusOK, &resp)
}
