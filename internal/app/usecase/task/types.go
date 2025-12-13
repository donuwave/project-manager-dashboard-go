package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskDTO struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
}

type CreateInput struct {
	Title       string
	Description *string
	Status      string // optional: "todo" | "in_progress" | "done"
}

type TasksRepository interface {
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]TaskDTO, error)
	CreateInProject(ctx context.Context, projectID uuid.UUID, in CreateInput) (TaskDTO, error)
}
