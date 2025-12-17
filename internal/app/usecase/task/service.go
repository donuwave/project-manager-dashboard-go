package task

import (
	"context"

	"github.com/google/uuid"
)

type TaskService interface {
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]TaskDTO, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (TaskDTO, error)
	Assign(ctx context.Context, taskID, actorID, userID uuid.UUID) error
	CreateInProject(ctx context.Context, projectID uuid.UUID, in CreateInput) (TaskDTO, error)
}
