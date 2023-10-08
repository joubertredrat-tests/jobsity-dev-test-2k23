package domain

import "context"

type TokenService interface {
	Generate(ctx context.Context, user User) (UserToken, error)
	Check(ctx context.Context, userToken UserToken) (User, error)
}
