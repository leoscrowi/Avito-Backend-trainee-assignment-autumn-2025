package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller interface {
	SetIsActive(w http.ResponseWriter, r *http.Request)
	GetReview(w http.ResponseWriter, r *http.Request)

	SetupRoutes(r chi.Router)
}
