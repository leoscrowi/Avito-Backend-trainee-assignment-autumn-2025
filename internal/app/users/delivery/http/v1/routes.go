package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/leoscrowi/pr-assignment-service/internal/config"
	"github.com/leoscrowi/pr-assignment-service/internal/middleware"
)

func (c *UsersController) SetupRoutes(r chi.Router, cfg *config.Config) {
	r.Route("/users", func(r chi.Router) {
		r.With(middleware.AuthMiddleware(cfg)).Get("/getReview/{user_id}", c.GetReview)
		r.With(middleware.AdminMiddleware(cfg)).Patch("/setIsActive", c.SetIsActive)
	})
}
