package v1

import "github.com/go-chi/chi/v5"

func (c *TeamsController) SetupRoutes(r chi.Router) {
	r.Route("/team", func(r chi.Router) {
		r.Get("/get", c.GetTeam)
		r.Post("/add", c.AddTeam)
	})
}
