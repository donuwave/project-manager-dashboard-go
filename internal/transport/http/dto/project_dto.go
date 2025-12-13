package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	UserID      string  `json:"userId"`
}

type ProjectUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type ProjectResponse struct {
	ID          uuid.UUID             `json:"id"`
	Name        string                `json:"name"`
	Description *string               `json:"description,omitempty"`
	CreatedAt   time.Time             `json:"createdAt"`
	Users       []ProjectUserResponse `json:"users"`
}

type InviteUserRequest struct {
	InviterID string `json:"inviterId"`
	UserID    string `json:"userId"`
}
