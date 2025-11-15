package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

type RouteSetup interface {
	SetupRoutes(r chi.Router)
}

type Server struct {
	Router      chi.Router
	Controllers []RouteSetup
}

func NewServer(db *sqlx.DB) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return &Server{
		Router:      r,
		Controllers: GetControllers(db),
	}
}

func (s *Server) SetupRoutes() {
	for _, controller := range s.Controllers {
		controller.SetupRoutes(s.Router)
	}
}
