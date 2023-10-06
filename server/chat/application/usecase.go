package application

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
)

type UsecaseUserRegister struct {
	userRepository domain.UserRepository
}

func NewUsecaseUserRegister(r domain.UserRepository) UsecaseUserRegister {
	return UsecaseUserRegister{
		userRepository: r,
	}
}

func (u UsecaseUserRegister) Execute(ctx context.Context, input UsecaseUserRegisterInput) (domain.User, error) {
	user, err := domain.NewUser("", input.Name, input.Email, input.Password)
	if err != nil {
		return domain.User{}, err
	}

	userGot, err := u.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	if userGot.Email == user.Email {
		return domain.User{}, NewErrUserAlreadyRegistered(user.Email)
	}

	return u.userRepository.Create(ctx, user)
}

type UsecaseUserLogin struct {
	userRepository domain.UserRepository
	tokenService   domain.TokenService
}

func NewUsecaseUserLogin(r domain.UserRepository, t domain.TokenService) UsecaseUserLogin {
	return UsecaseUserLogin{
		userRepository: r,
		tokenService:   t,
	}
}

func (u UsecaseUserLogin) Execute(ctx context.Context, input UsecaseUserLoginInput) (domain.UserToken, error) {
	user, err := domain.NewUser("", "", input.Email, input.Password)
	if err != nil {
		return domain.UserToken{}, err
	}
	userGot, err := u.userRepository.GetAuthenticated(ctx, user)
	if err != nil {
		return domain.UserToken{}, err
	}

	return u.tokenService.Generate(ctx, userGot)
}
