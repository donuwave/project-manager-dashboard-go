package project

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateInput struct {
	Name        string
	Description *string
	OwnerID     uuid.UUID
}

type ProjectMemberDTO struct {
	UserID uuid.UUID
	Name   string
	Email  string
	Role   string
}

type ProjectTaskDTO struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
}

type ProjectDTO struct {
	ID          uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time

	Members []ProjectMemberDTO
	Tasks   []ProjectTaskDTO
}

type ProjectRepository interface {
	Create(ctx context.Context, in CreateInput) (ProjectDTO, error)
	GetByID(ctx context.Context, id uuid.UUID) (ProjectDTO, error)
	List(ctx context.Context, limit, offset int) ([]ProjectDTO, error)
	ProjectExists(ctx context.Context, projectID uuid.UUID) (bool, error)
	UserExists(ctx context.Context, userID uuid.UUID) (bool, error)
	IsMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error)
	AddMember(ctx context.Context, projectID, userID uuid.UUID, role string) error
}
