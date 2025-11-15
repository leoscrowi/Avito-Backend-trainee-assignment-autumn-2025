package teams

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Controller interface {
	AddTeam(w http.ResponseWriter, r *http.Request)
	GetTeam(w http.ResponseWriter, r *http.Request)

	SetupRoutes(r chi.Router)
}
