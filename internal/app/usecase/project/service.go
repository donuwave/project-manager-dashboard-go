package project

import (
	"context"

	"github.com/google/uuid"
)

type ProjectService interface {
	Create(ctx context.Context, in CreateInput) (ProjectDTO, error)
	GetByID(ctx context.Context, id uuid.UUID) (ProjectDTO, error)
	List(ctx context.Context, limit, offset int) ([]ProjectDTO, error)

	Invite(ctx context.Context, projectID, inviterID, userID uuid.UUID) error
}
