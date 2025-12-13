package user

import (
	"context"

	"github.com/google/uuid"
)

func NewUserUsecase(repo UserRepository) UserRepository {
	return &userUC{repo: repo}
}

func (u *userUC) List(ctx context.Context, limit int) ([]User, error) {
	return u.repo.List(ctx, limit)
}

func (u *userUC) Create(ctx context.Context, in CreateUserInput) (User, error) {
	return u.repo.Create(ctx, in)
}

func (u *userUC) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	return u.repo.GetByID(ctx, id)
}
