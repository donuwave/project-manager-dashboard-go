package dto

import (
	"time"

	"github.com/google/uuid"
)

type DeleteTaskRequest struct {
	ActorID string `json:"actorId"`
}

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status,omitempty"`
}

type TaskAssigneeResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type TaskResponse struct {
	ID          uuid.UUID             `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description,omitempty"`
	Status      string                `json:"status"`
	CreatedAt   time.Time             `json:"createdAt"`
	Assignee    *TaskAssigneeResponse `json:"assignees"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"` // "todo" | "in_progress" | "done"
}

type AssignTaskRequest struct {
	ActorID string `json:"actorId"`          // кто делает действие
	UserID  string `json:"userId,omitempty"` // на кого назначаем (если пусто - на себя)
}
