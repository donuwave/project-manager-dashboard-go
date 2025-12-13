package http

import (
	"encoding/json"
	"net/http"
	"project-manager-dashboard-go/internal/app/usecase/user"
	"project-manager-dashboard-go/internal/transport/http/dto"
	"project-manager-dashboard-go/internal/transport/http/mapper"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	uc user.UserRepository
}

func NewUserHandler(uc user.UserRepository) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.List(r.Context(), 50)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	out := make([]dto.UserDTO, 0, len(users))
	for _, u := range users {
		out = append(out, mapper.ToUserDTO(u))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Name == "" {
		http.Error(w, "email and name required", http.StatusBadRequest)
		return
	}

	u, err := h.uc.Create(r.Context(), user.CreateUserInput{
		Email:   req.Email,
		Name:    req.Name,
		Country: req.Country,
	})
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(mapper.ToUserDTO(u))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	u, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(mapper.ToUserDTO(u))
}
