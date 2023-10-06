package domain

import "context"

type TokenService interface {
	Generate(ctx context.Context, user User) (UserToken, error)
}
