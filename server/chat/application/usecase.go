package application

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
)

type UsecaseUserRegister struct {
	userRepository domain.UserRepository
}

func NewUsecaseUserRegister(userRepository domain.UserRepository) UsecaseUserRegister {
	return UsecaseUserRegister{
		userRepository: userRepository,
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

	return u.userRepository.Persist(ctx, user)
}
