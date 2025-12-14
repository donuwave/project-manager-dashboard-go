package task

import (
	"context"
	"errors"
	"project-manager-dashboard-go/ent/task"

	"github.com/google/uuid"
	"project-manager-dashboard-go/ent"
	"project-manager-dashboard-go/ent/project"
	"project-manager-dashboard-go/ent/projecttask"
)

type EntRepo struct{ client *ent.Client }

func NewEntRepo(c *ent.Client) *EntRepo { return &EntRepo{client: c} }

func (r *EntRepo) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]TaskDTO, error) {
	rows, err := r.client.ProjectTask.
		Query().
		Where(projecttask.HasProjectWith(project.IDEQ(projectID))).
		WithTask(). // подтянуть сам Task
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]TaskDTO, 0, len(rows))
	for _, row := range rows {
		t := row.Edges.Task
		if t == nil {
			continue
		}
		out = append(out, TaskDTO{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      string(t.Status),
			CreatedAt:   t.CreatedAt,
		})
	}
	return out, nil
}

func (r *EntRepo) CreateInProject(ctx context.Context, projectID uuid.UUID, in CreateInput) (TaskDTO, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer func() { _ = tx.Rollback() }()

	tc := tx.Task.Create().SetTitle(in.Title)
	if in.Description != nil {
		tc.SetDescription(*in.Description)
	}
	if in.Status != "" {
		tc.SetStatus(task.Status(in.Status))
	}

	t, err := tc.Save(ctx)
	if err != nil {
		return TaskDTO{}, err
	}

	_, err = tx.ProjectTask.
		Create().
		SetProjectID(projectID).
		SetTaskID(t.ID).
		Save(ctx)
	if err != nil {
		return TaskDTO{}, err
	}

	if err := tx.Commit(); err != nil {
		return TaskDTO{}, err
	}

	return TaskDTO{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt,
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (TaskDTO, error) {
	u := r.client.Task.UpdateOneID(id)

	if in.Title != nil {
		u.SetTitle(*in.Title)
	}
	if in.Description != nil {
		if *in.Description == "" {
			u.ClearDescription()
		} else {
			u.SetDescription(*in.Description)
		}
	}
	if in.Status != nil {
		u.SetStatus(task.Status(*in.Status))
	}

	t, err := u.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return TaskDTO{}, errors.New("task not found")
		}
		return TaskDTO{}, err
	}

	return TaskDTO{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt,
	}, nil
}
