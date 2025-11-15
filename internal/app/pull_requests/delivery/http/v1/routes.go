package v1

import "github.com/go-chi/chi/v5"

func (c *PullRequestController) SetupRoutes(r chi.Router) {
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", c.CreatePullRequest)
		r.Patch("/reassign", c.ReassignPullRequest)
		r.Patch("/merge", c.MergePullRequest)
	})
}
