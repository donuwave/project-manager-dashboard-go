package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskAssigneeDTO struct {
	UserID uuid.UUID
	Name   string
	Email  string
}

type TaskDTO struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time

	Assignee *TaskAssigneeDTO
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
	GetProjectIDByTask(ctx context.Context, taskID uuid.UUID) (uuid.UUID, error)
	UserExists(ctx context.Context, userID uuid.UUID) (bool, error)
	IsProjectMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error)
	GetMemberRole(ctx context.Context, projectID, userID uuid.UUID) (string, error)
	GetAssignee(ctx context.Context, taskID uuid.UUID) (*TaskAssigneeDTO, error)
	SetAssignee(ctx context.Context, taskID, userID uuid.UUID) error
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
}
