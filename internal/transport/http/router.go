package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(userH *UserHandler, projectH *ProjectHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	r.Get("/users", userH.ListUsers)
	r.Post("/users", userH.CreateUser)
	r.Get("/users/{id}", userH.GetUser)

	r.Get("/projects", projectH.ListProjects)
	r.Post("/projects", projectH.CreateProject)
	r.Get("/projects/{id}", projectH.GetProject)
	r.Post("/projects/{id}/invite", projectH.Invite)

	return r
}
