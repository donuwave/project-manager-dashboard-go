package user

import (
	"context"
	"github.com/google/uuid"
)

type User struct {
	ID      uuid.UUID
	Email   string
	Name    string
	Country string
}

type CreateUserInput struct {
	Email   string
	Name    string
	Country string
}

type UserRepository interface {
	List(ctx context.Context, limit int) ([]User, error)
	Create(ctx context.Context, in CreateUserInput) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
}

type userUC struct {
	repo UserRepository
}
