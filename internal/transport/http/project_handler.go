package http

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"project-manager-dashboard-go/internal/app/usecase/project"
	"project-manager-dashboard-go/internal/transport/http/dto"
)

type ProjectHandler struct {
	uc project.ProjectService
}

func NewProjectHandler(uc project.ProjectService) *ProjectHandler {
	return &ProjectHandler{uc: uc}
}

func (h *ProjectHandler) ListProjects(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()

	limit := 50
	offset := 0

	if s := r.URL.Query().Get("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			limit = v
		}
	}
	if s := r.URL.Query().Get("offset"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			offset = v
		}
	}

	items, err := h.uc.List(ctx, limit, offset)
	if err != nil {
		writeJSON(w, stdhttp.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	out := make([]dto.ProjectResponse, 0, len(items))
	for _, p := range items {
		out = append(out, dto.ProjectResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
		})
	}
	writeJSON(w, stdhttp.StatusOK, out)
}

func (h *ProjectHandler) CreateProject(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()

	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	ownerID, err := uuid.Parse(req.UserID)
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid userId"})
		return
	}

	created, err := h.uc.Create(ctx, project.CreateInput{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
	})
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, stdhttp.StatusCreated, dto.ProjectResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		CreatedAt:   created.CreatedAt,
	})
}

func (h *ProjectHandler) GetProject(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	p, err := h.uc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, project.ErrNotFound) {
			writeJSON(w, stdhttp.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, stdhttp.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	users := make([]dto.ProjectUserResponse, 0, len(p.Members))
	for _, m := range p.Members {
		users = append(users, dto.ProjectUserResponse{
			ID:    m.UserID,
			Name:  m.Name,
			Email: m.Email,
			Role:  m.Role,
		})
	}

	tasks := make([]dto.TaskResponse, 0, len(p.Tasks))
	for _, t := range p.Tasks {
		tasks = append(tasks, dto.TaskResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			CreatedAt:   t.CreatedAt,
		})
	}

	writeJSON(w, stdhttp.StatusOK, dto.ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		Users:       users,
		Tasks:       tasks,
	})
}

func (h *ProjectHandler) Invite(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid project id"})
		return
	}

	var req dto.InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	inviterID, err := uuid.Parse(req.InviterID)
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid inviterId"})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid userId"})
		return
	}

	err = h.uc.Invite(ctx, projectID, inviterID, userID)
	if err != nil {
		switch {
		case errors.Is(err, project.ErrNotFound):
			writeJSON(w, stdhttp.StatusNotFound, map[string]string{"error": "project not found"})
		case errors.Is(err, project.ErrUserNotFound):
			writeJSON(w, stdhttp.StatusNotFound, map[string]string{"error": "user not found"})
		case errors.Is(err, project.ErrForbidden):
			writeJSON(w, stdhttp.StatusForbidden, map[string]string{"error": "forbidden"})
		case errors.Is(err, project.ErrAlreadyMember):
			writeJSON(w, stdhttp.StatusConflict, map[string]string{"error": "already member"})
		default:
			writeJSON(w, stdhttp.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	writeJSON(w, stdhttp.StatusCreated, map[string]string{"status": "invited"})
}

func writeJSON(w stdhttp.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
