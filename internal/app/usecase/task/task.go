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

func (uc *UseCase) CreateInProject(ctx context.Context, projectID uuid.UUID, in CreateInput) (TaskDTO, error) {
	if strings.TrimSpace(in.Title) == "" {
		return TaskDTO{}, errors.New("title is required")
	}
	// status можно валидировать, но можно и пропускать (пусть ent default поставит)
	return uc.repo.CreateInProject(ctx, projectID, in)
}
