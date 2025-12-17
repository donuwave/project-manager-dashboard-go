package task

import (
	"context"
	"errors"
	"project-manager-dashboard-go/ent/projectuser"
	"project-manager-dashboard-go/ent/task"
	"project-manager-dashboard-go/ent/user"

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
		WithTask(func(tq *ent.TaskQuery) {
			tq.WithAssignee()
		}).
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

		var assignee *TaskAssigneeDTO
		if u := t.Edges.Assignee; u != nil {
			assignee = &TaskAssigneeDTO{
				UserID: u.ID,
				Name:   u.Name,
				Email:  u.Email,
			}
		}

		out = append(out, TaskDTO{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      string(t.Status),
			CreatedAt:   t.CreatedAt,

			Assignee: assignee,
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

func (r *EntRepo) GetProjectIDByTask(ctx context.Context, taskID uuid.UUID) (uuid.UUID, error) {
	pt, err := r.client.ProjectTask.
		Query().
		Where(projecttask.HasTaskWith(task.IDEQ(taskID))).
		WithProject().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	if pt.Edges.Project == nil {
		return uuid.Nil, nil
	}
	return pt.Edges.Project.ID, nil
}

func (r *EntRepo) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	return r.client.User.Query().Where(user.IDEQ(userID)).Exist(ctx)
}

func (r *EntRepo) IsProjectMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	return r.client.ProjectUser.
		Query().
		Where(
			projectuser.HasProjectWith(project.IDEQ(projectID)),
			projectuser.HasUserWith(user.IDEQ(userID)),
		).
		Exist(ctx)
}

func (r *EntRepo) GetMemberRole(ctx context.Context, projectID, userID uuid.UUID) (string, error) {
	m, err := r.client.ProjectUser.
		Query().
		Where(
			projectuser.HasProjectWith(project.IDEQ(projectID)),
			projectuser.HasUserWith(user.IDEQ(userID)),
		).
		Only(ctx)
	if err != nil {
		return "", err
	}
	return string(m.Role), nil
}

func (r *EntRepo) GetAssignee(ctx context.Context, taskID uuid.UUID) (*TaskAssigneeDTO, error) {
	t, err := r.client.Task.
		Query().
		Where(task.IDEQ(taskID)).
		WithAssignee().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	u := t.Edges.Assignee
	if u == nil {
		return nil, nil
	}

	return &TaskAssigneeDTO{
		UserID: u.ID,
		Name:   u.Name,
		Email:  u.Email,
	}, nil
}

func (r *EntRepo) SetAssignee(ctx context.Context, taskID, userID uuid.UUID) error {
	return r.client.Task.
		UpdateOneID(taskID).
		SetAssigneeID(userID).
		Exec(ctx)
}

func (r *EntRepo) ClearAssignee(ctx context.Context, taskID uuid.UUID) error {
	return r.client.Task.
		UpdateOneID(taskID).
		ClearAssigneeID().
		Exec(ctx)
}
