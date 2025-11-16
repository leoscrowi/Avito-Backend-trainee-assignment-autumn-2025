package pull_requests

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leoscrowi/pr-assignment-service/internal/config"
)

type Controller interface {
	CreatePullRequest(w http.ResponseWriter, r *http.Request)
	MergePullRequest(w http.ResponseWriter, r *http.Request)
	ReassignPullRequest(w http.ResponseWriter, r *http.Request)

	SetupRoutes(r chi.Router, cfg *config.Config)
}
