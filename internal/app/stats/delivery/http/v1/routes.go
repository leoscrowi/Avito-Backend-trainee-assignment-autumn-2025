package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/leoscrowi/pr-assignment-service/internal/config"
	"github.com/leoscrowi/pr-assignment-service/internal/middleware"
)

func (c *StatsController) SetupRoutes(r chi.Router, cfg *config.Config) {
	r.Route("/stats", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg))
		r.Get("/users", c.GetPullRequestStats)
	})

}
