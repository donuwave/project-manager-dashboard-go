package user

import (
	"context"

	"github.com/google/uuid"
	"project-manager-dashboard-go/ent"
	"project-manager-dashboard-go/ent/user"
)

type UserRepo struct {
	ent *ent.Client
}

func NewUserRepo(entClient *ent.Client) *UserRepo {
	return &UserRepo{ent: entClient}
}

func (r *UserRepo) List(ctx context.Context, limit int) ([]User, error) {
	rows, err := r.ent.User.Query().Limit(limit).All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]User, 0, len(rows))
	for _, u := range rows {
		out = append(out, User{
			ID:      u.ID,
			Email:   u.Email,
			Name:    u.Name,
			Country: u.Country,
		})
	}
	return out, nil
}

func (r *UserRepo) Create(ctx context.Context, in CreateUserInput) (User, error) {
	q := r.ent.User.Create().SetEmail(in.Email).SetName(in.Name)
	if in.Country != "" {
		q.SetCountry(in.Country)
	}

	u, err := q.Save(ctx)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:      u.ID,
		Email:   u.Email,
		Name:    u.Name,
		Country: u.Country,
	}, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	u, err := r.ent.User.Query().Where(user.IDEQ(id)).Only(ctx)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:      u.ID,
		Email:   u.Email,
		Name:    u.Name,
		Country: u.Country,
	}, nil
}
