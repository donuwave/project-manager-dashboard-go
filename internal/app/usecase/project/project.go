package project

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type UseCase struct {
	repo ProjectRepository
}

func NewProjectUsecase(repo ProjectRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(ctx context.Context, in CreateInput) (ProjectDTO, error) {
	if strings.TrimSpace(in.Name) == "" {
		return ProjectDTO{}, errors.New("name is required")
	}
	return uc.repo.Create(ctx, in)
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (ProjectDTO, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (ProjectDTO, error) {
	if in.Name != nil && strings.TrimSpace(*in.Name) == "" {
		return ProjectDTO{}, errors.New("name cannot be empty")
	}
	return uc.repo.Update(ctx, id, in)
}

func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]ProjectDTO, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.List(ctx, limit, offset)
}

func (uc *UseCase) Invite(ctx context.Context, projectID, inviterID, userID uuid.UUID) error {
	ok, err := uc.repo.ProjectExists(ctx, projectID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}

	ok, err = uc.repo.UserExists(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	isInviterMember, err := uc.repo.IsMember(ctx, projectID, inviterID)
	if err != nil {
		return err
	}
	if !isInviterMember {
		return ErrForbidden
	}

	isAlready, err := uc.repo.IsMember(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if isAlready {
		return ErrAlreadyMember
	}

	return uc.repo.AddMember(ctx, projectID, userID, "member")
}
