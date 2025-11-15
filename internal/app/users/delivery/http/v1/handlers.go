package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/leoscrowi/pr-assignment-service/internal/app/users"
	"github.com/leoscrowi/pr-assignment-service/internal/app/users/dtos"
	"github.com/leoscrowi/pr-assignment-service/internal/utils"
)

type UsersController struct {
	usecase users.Usecase
}

func NewUsersController(usecase users.Usecase) *UsersController {
	return &UsersController{usecase: usecase}
}

func (c *UsersController) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req dtos.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", err))
		return
	}

	if req.UserID == "" {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", fmt.Errorf("wrong json format")))
		return
	}

	user, err := c.usecase.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		domain.WriteError(w, domain.ConvertToErrorResponse(err))
		return
	}

	var resp = dtos.SetIsActiveResponse{
		User: user,
	}
	utils.WriteHeader(w, http.StatusOK, &resp)
}

func (c *UsersController) GetReview(w http.ResponseWriter, r *http.Request) {
	var req dtos.GetReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", err))
		return
	}

	if req.UserID == "" {
		domain.WriteError(w, domain.NewError(domain.BAD_REQUEST, "bad request", fmt.Errorf("wrong json format")))
		return
	}

	prs, err := c.usecase.GetReview(r.Context(), req.UserID)
	if err != nil {
		domain.WriteError(w, domain.ConvertToErrorResponse(err))
	}

	var resp = dtos.GetReviewResponse{UserID: req.UserID, PullRequests: prs}
	utils.WriteHeader(w, http.StatusOK, &resp)
}
