package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status,omitempty"`
}

type TaskResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"` // "todo" | "in_progress" | "done"
}
