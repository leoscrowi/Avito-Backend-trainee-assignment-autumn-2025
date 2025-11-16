package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/leoscrowi/pr-assignment-service/internal/config"
	"github.com/leoscrowi/pr-assignment-service/internal/middleware"
)

func (c *TeamsController) SetupRoutes(r chi.Router, cfg *config.Config) {
	r.Route("/team", func(r chi.Router) {
		r.With(middleware.AuthMiddleware(cfg)).Get("/get/{team_name}", c.GetTeam)
		r.Post("/add", c.AddTeam)
	})
}
