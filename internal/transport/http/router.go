package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(userH *UserHandler, projectH *ProjectHandler, taskH *TaskHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	r.Get("/users", userH.ListUsers)
	r.Post("/users", userH.CreateUser)
	r.Get("/users/{id}", userH.GetUser)

	r.Get("/projects", projectH.ListProjects)
	r.Post("/projects", projectH.CreateProject)
	r.Get("/projects/{id}", projectH.GetProject)
	r.Patch("/projects/{id}", projectH.UpdateProject)
	r.Post("/projects/{id}/invite", projectH.Invite)
	r.Get("/projects/{id}/tasks", taskH.ListByProject)
	r.Post("/projects/{id}/tasks", taskH.CreateInProject)
	r.Delete("/projects/{id}", projectH.DeleteProject)

	r.Patch("/tasks/{id}", taskH.UpdateTask)
	r.Post("/tasks/{id}/assign", taskH.Assign)
	r.Delete("/tasks/{id}", taskH.DeleteTask)

	return r
}
