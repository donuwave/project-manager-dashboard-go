package task

import (
	"context"

	"github.com/google/uuid"
)

type TaskService interface {
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]TaskDTO, error)
	CreateInProject(ctx context.Context, projectID uuid.UUID, in CreateInput) (TaskDTO, error)
}
