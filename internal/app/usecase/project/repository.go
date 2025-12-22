package project

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"project-manager-dashboard-go/ent/project"
	"project-manager-dashboard-go/ent/projecttask"
	"project-manager-dashboard-go/ent/projectuser"
	"project-manager-dashboard-go/ent/task"
	"project-manager-dashboard-go/ent/user"

	"github.com/google/uuid"
	"project-manager-dashboard-go/ent"
)

type EntRepo struct {
	client *ent.Client
}

func NewEntRepo(c *ent.Client) *EntRepo {
	return &EntRepo{client: c}
}

func (r *EntRepo) Create(ctx context.Context, in CreateInput) (ProjectDTO, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return ProjectDTO{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	pc := tx.Project.Create().SetName(in.Name)
	if in.Description != nil {
		pc.SetDescription(*in.Description)
	}

	p, err := pc.Save(ctx)
	if err != nil {
		return ProjectDTO{}, err
	}

	_, err = tx.ProjectUser.
		Create().
		SetProjectID(p.ID).
		SetUserID(in.OwnerID).
		SetRole("owner").
		Save(ctx)
	if err != nil {
		return ProjectDTO{}, err
	}

	if err = tx.Commit(); err != nil {
		return ProjectDTO{}, err
	}

	return ProjectDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}, nil
}

func (r *EntRepo) GetByID(ctx context.Context, id uuid.UUID) (ProjectDTO, error) {
	p, err := r.client.Project.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return ProjectDTO{}, ErrNotFound
		}
		return ProjectDTO{}, err
	}

	memberships, err := p.QueryMemberships().
		WithUser().
		All(ctx)
	if err != nil {
		return ProjectDTO{}, err
	}

	members := make([]ProjectMemberDTO, 0, len(memberships))
	for _, m := range memberships {
		u := m.Edges.User
		if u == nil {
			continue
		}
		members = append(members, ProjectMemberDTO{
			UserID: u.ID,
			Name:   u.Name,
			Email:  u.Email,
			Role:   string(m.Role),
		})
	}

	projectTasks, err := p.QueryProjectTasks().
		WithTask().
		Order(ent.Asc(projecttask.FieldPosition), ent.Asc(projecttask.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return ProjectDTO{}, err
	}

	tasks := make([]ProjectTaskDTO, 0, len(projectTasks))
	for _, pt := range projectTasks {
		t := pt.Edges.Task
		if t == nil {
			continue
		}
		tasks = append(tasks, ProjectTaskDTO{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      string(t.Status),
			CreatedAt:   t.CreatedAt,
			Position:    pt.Position,
		})
	}

	return ProjectDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		Members:     members,
		Tasks:       tasks,
	}, nil
}

func (r *EntRepo) List(ctx context.Context, limit, offset int) ([]ProjectDTO, error) {
	items, err := r.client.Project.
		Query().
		Order(project.ByCreatedAt(sql.OrderDesc())).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ProjectDTO, 0, len(items))
	for _, p := range items {
		out = append(out, ProjectDTO{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
		})
	}
	return out, nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (ProjectDTO, error) {
	pu := r.client.Project.UpdateOneID(id)

	if in.Name != nil {
		pu.SetName(*in.Name)
	}
	if in.Description != nil {
		if *in.Description == "" {
			pu.ClearDescription()
		} else {
			pu.SetDescription(*in.Description)
		}
	}

	p, err := pu.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ProjectDTO{}, ErrNotFound
		}
		return ProjectDTO{}, err
	}

	return ProjectDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}, nil
}

func (r *EntRepo) ProjectExists(ctx context.Context, projectID uuid.UUID) (bool, error) {
	return r.client.Project.
		Query().
		Where(project.IDEQ(projectID)).
		Exist(ctx)
}

func (r *EntRepo) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	return r.client.User.
		Query().
		Where(user.IDEQ(userID)).
		Exist(ctx)
}

func (r *EntRepo) IsMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	return r.client.ProjectUser.
		Query().
		Where(
			projectuser.HasProjectWith(project.IDEQ(projectID)),
			projectuser.HasUserWith(user.IDEQ(userID)),
		).
		Exist(ctx)
}

func (r *EntRepo) AddMember(ctx context.Context, projectID, userID uuid.UUID, role string) error {
	_, err := r.client.ProjectUser.
		Create().
		SetProjectID(projectID).
		SetUserID(userID).
		SetRole(projectuser.Role(role)).
		Save(ctx)

	if err != nil && ent.IsConstraintError(err) {
		return ErrAlreadyMember
	}
	return err
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
func (r *EntRepo) DeleteProject(ctx context.Context, projectID uuid.UUID) (err error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	taskIDs, err := tx.ProjectTask.
		Query().
		Where(projecttask.HasProjectWith(project.IDEQ(projectID))).
		QueryTask().
		IDs(ctx)
	if err != nil {
		return err
	}

	_, err = tx.ProjectTask.
		Delete().
		Where(projecttask.HasProjectWith(project.IDEQ(projectID))).
		Exec(ctx)
	if err != nil {
		return err
	}

	if len(taskIDs) > 0 {
		_, err = tx.Task.
			Delete().
			Where(task.IDIn(taskIDs...)).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	_, err = tx.ProjectUser.
		Delete().
		Where(projectuser.HasProjectWith(project.IDEQ(projectID))).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Project.
		Delete().
		Where(project.IDEQ(projectID)).
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
