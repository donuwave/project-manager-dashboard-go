package http

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"
	"project-manager-dashboard-go/internal/transport/http/dto"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"project-manager-dashboard-go/internal/app/usecase/task"
)

type TaskHandler struct {
	uc task.TaskService
}

func NewTaskHandler(uc task.TaskService) *TaskHandler {
	return &TaskHandler{uc: uc}
}

func (h *TaskHandler) ListByProject(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid project id"})
		return
	}

	limit, offset := 50, 0
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

	items, err := h.uc.ListByProject(ctx, projectID, limit, offset)
	if err != nil {
		writeJSON(w, stdhttp.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, stdhttp.StatusOK, items)
}

func (h *TaskHandler) CreateInProject(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid project id"})
		return
	}

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	created, err := h.uc.CreateInProject(ctx, projectID, task.CreateInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, stdhttp.StatusCreated, dto.TaskResponse{
		ID:          created.ID,
		Title:       created.Title,
		Description: created.Description,
		Status:      created.Status,
		CreatedAt:   created.CreatedAt,
	})
}

func (h *TaskHandler) UpdateTask(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	updated, err := h.uc.Update(ctx, id, task.UpdateInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		if errors.Is(err, errors.New("task not found")) {
			writeJSON(w, stdhttp.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, stdhttp.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, stdhttp.StatusOK, updated)
}
