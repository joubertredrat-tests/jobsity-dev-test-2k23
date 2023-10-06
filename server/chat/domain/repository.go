package domain

import "context"

type UserRepository interface {
	Persist(ctx context.Context, user User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}
