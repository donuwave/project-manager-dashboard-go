package task

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type UseCase struct{ repo TasksRepository }

func NewTasksUseCase(repo TasksRepository) *UseCase { return &UseCase{repo: repo} }

func (uc *UseCase) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]TaskDTO, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.ListByProject(ctx, projectID, limit, offset)
}

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (TaskDTO, error) {
	if in.Title != nil && strings.TrimSpace(*in.Title) == "" {
		return TaskDTO{}, errors.New("title cannot be empty")
	}
	if in.Status != nil {
		switch *in.Status {
		case "todo", "in_progress", "done":
		default:
			return TaskDTO{}, errors.New("invalid status")
		}
	}
	if in.Position != nil && *in.Position < 0 {
		return TaskDTO{}, errors.New("position must be >= 0")
	}

	return uc.repo.Update(ctx, id, in)
}

func (uc *UseCase) CreateInProject(ctx context.Context, projectID uuid.UUID, in CreateInput) (TaskDTO, error) {
	if strings.TrimSpace(in.Title) == "" {
		return TaskDTO{}, errors.New("title is required")
	}
	return uc.repo.CreateInProject(ctx, projectID, in)
}

func (uc *UseCase) Assign(ctx context.Context, taskID, actorID, userID uuid.UUID) error {
	projectID, err := uc.repo.GetProjectIDByTask(ctx, taskID)
	if err != nil {
		return err
	}
	if projectID == uuid.Nil {
		return errors.New("task not found")
	}

	ok, err := uc.repo.UserExists(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("user not found")
	}

	isActorMember, err := uc.repo.IsProjectMember(ctx, projectID, actorID)
	if err != nil {
		return err
	}
	if !isActorMember {
		return errors.New("forbidden")
	}

	isTargetMember, err := uc.repo.IsProjectMember(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if !isTargetMember {
		return errors.New("user not in project")
	}

	if actorID != userID {
		role, err := uc.repo.GetMemberRole(ctx, projectID, actorID)
		if err != nil {
			return err
		}
		if role != "owner" {
			return errors.New("forbidden")
		}
	}

	cur, err := uc.repo.GetAssignee(ctx, taskID)
	if err != nil {
		return err
	}
	if cur != nil && cur.UserID == userID {
		return errors.New("already assigned")
	}

	return uc.repo.SetAssignee(ctx, taskID, userID)
}

func (uc *UseCase) Delete(ctx context.Context, taskID, actorID uuid.UUID) error {
	projectID, err := uc.repo.GetProjectIDByTask(ctx, taskID)
	if err != nil {
		return err
	}

	role, err := uc.repo.GetMemberRole(ctx, projectID, actorID)
	if err != nil {
		return err
	}
	if role != "owner" {
		return ErrForbidden
	}

	return uc.repo.DeleteTask(ctx, taskID)
}
