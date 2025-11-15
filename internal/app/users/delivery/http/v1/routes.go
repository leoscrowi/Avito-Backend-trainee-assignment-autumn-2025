package v1

import "github.com/go-chi/chi/v5"

func (c *UsersController) SetupRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/getReview", c.GetReview)
		r.Patch("/setIsActive", c.SetIsActive)
	})
}
