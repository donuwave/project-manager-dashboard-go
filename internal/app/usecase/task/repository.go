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
	enttask "project-manager-dashboard-go/ent/task"
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
		Order(ent.Asc(projecttask.FieldPosition), ent.Asc(projecttask.FieldCreatedAt)).
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
			Position:    row.Position,

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

	maxPos := 0
	last, err := tx.ProjectTask.
		Query().
		Where(projecttask.HasProjectWith(project.IDEQ(projectID))).
		Order(ent.Desc(projecttask.FieldPosition)).
		First(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return TaskDTO{}, err
		}
	} else {
		maxPos = last.Position + 1
	}

	pt, err := tx.ProjectTask.
		Create().
		SetProjectID(projectID).
		SetTaskID(t.ID).
		SetPosition(maxPos + 1).
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
		Position:    pt.Position, // если добавил поле
	}, nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (TaskDTO, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return TaskDTO{}, err
	}
	defer func() { _ = tx.Rollback() }()

	u := tx.Task.UpdateOneID(id)

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

	pt, err := tx.ProjectTask.
		Query().
		Where(projecttask.HasTaskWith(enttask.IDEQ(id))).
		WithProject().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return TaskDTO{}, errors.New("task not found")
		}
		var nse *ent.NotSingularError
		if errors.As(err, &nse) {
			return TaskDTO{}, errors.New("task not found")
		}
		return TaskDTO{}, err
	}

	if in.Position != nil {
		target := *in.Position
		if target < 0 {
			target = 0
		}

		if pt.Edges.Project == nil {
			return TaskDTO{}, errors.New("task not found")
		}
		projectID := pt.Edges.Project.ID
		curPos := pt.Position

		count, err := tx.ProjectTask.
			Query().
			Where(projecttask.HasProjectWith(project.IDEQ(projectID))).
			Count(ctx)
		if err != nil {
			return TaskDTO{}, err
		}
		if count <= 0 {
			target = 0
		} else if target > count-1 {
			target = count - 1
		}

		if target != curPos {
			tempPos := count

			_, err = tx.ProjectTask.
				UpdateOneID(pt.ID).
				SetPosition(tempPos).
				Save(ctx)
			if err != nil {
				return TaskDTO{}, err
			}

			if target < curPos {
				_, err = tx.ProjectTask.
					Update().
					Where(
						projecttask.HasProjectWith(project.IDEQ(projectID)),
						projecttask.PositionGTE(target),
						projecttask.PositionLT(curPos),
						projecttask.IDNEQ(pt.ID),
					).
					AddPosition(1).
					Save(ctx)
				if err != nil {
					return TaskDTO{}, err
				}
			} else {
				_, err = tx.ProjectTask.
					Update().
					Where(
						projecttask.HasProjectWith(project.IDEQ(projectID)),
						projecttask.PositionGT(curPos),
						projecttask.PositionLTE(target),
						projecttask.IDNEQ(pt.ID),
					).
					AddPosition(-1).
					Save(ctx)
				if err != nil {
					return TaskDTO{}, err
				}
			}

			pt, err = tx.ProjectTask.
				UpdateOneID(pt.ID).
				SetPosition(target).
				Save(ctx)
			if err != nil {
				return TaskDTO{}, err
			}
		}
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
		Position:    pt.Position,
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
			return uuid.Nil, ErrNotFound
		}
		return uuid.Nil, err
	}
	if pt.Edges.Project == nil {
		return uuid.Nil, ErrNotFound
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
		if ent.IsNotFound(err) {
			return "", ErrForbidden
		}
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

func (r *EntRepo) DeleteTask(ctx context.Context, taskID uuid.UUID) (err error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ProjectTask.
		Delete().
		Where(projecttask.HasTaskWith(enttask.IDEQ(taskID))).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = tx.Task.DeleteOneID(taskID).Exec(ctx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
