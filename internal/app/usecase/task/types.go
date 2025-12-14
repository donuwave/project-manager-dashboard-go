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

type UpdateInput struct {
	Title       *string
	Description *string
	Status      *string
}

type CreateInput struct {
	Title       string
	Description *string
	Status      string
}

type TasksRepository interface {
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]TaskDTO, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (TaskDTO, error)
	CreateInProject(ctx context.Context, projectID uuid.UUID, in CreateInput) (TaskDTO, error)
}
