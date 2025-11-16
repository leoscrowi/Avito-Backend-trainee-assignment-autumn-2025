package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/leoscrowi/pr-assignment-service/internal/config"
	"github.com/leoscrowi/pr-assignment-service/internal/middleware"
)

func (c *PullRequestController) SetupRoutes(r chi.Router, cfg *config.Config) {
	r.Route("/pullRequest", func(r chi.Router) {
		r.Use(middleware.AdminMiddleware(cfg))
		r.Post("/create", c.CreatePullRequest)
		r.Patch("/reassign", c.ReassignPullRequest)
		r.Patch("/merge", c.MergePullRequest)
	})
}
