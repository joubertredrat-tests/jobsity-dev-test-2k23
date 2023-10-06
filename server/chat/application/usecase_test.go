package application_test

import (
	"context"
	"errors"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/application"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"joubertredrat-tests/jobsity-dev-test-2k23/pkg/chat/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUsecaseUserRegister(t *testing.T) {
	tests := []struct {
		name                     string
		userRepositoryDependency func(ctrl *gomock.Controller) domain.UserRepository
		input                    application.UsecaseUserRegisterInput
		userExpected             domain.User
		errExpected              error
	}{
		{
			name: "Test register user with success",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := mock.NewMockUserRepository(ctrl)
				repo.
					EXPECT().
					GetByEmail(gomock.Any(), gomock.Eq("foo@bar.tld")).
					Return(domain.User{}, nil).
					Times(1)
				repo.
					EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
					Return(domain.User{
						ID:       "ID",
						Name:     "Sr Foo",
						Email:    "foo@bar.tld",
						Password: "password-hashed",
					}, nil).
					Times(1)

				return repo
			},
			input: application.UsecaseUserRegisterInput{
				Name:     "Sr Foo",
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userExpected: domain.User{
				ID:       "ID",
				Name:     "Sr Foo",
				Email:    "foo@bar.tld",
				Password: "password-hashed",
			},
			errExpected: nil,
		},
		{
			name: "Test register user with invalid e-mail",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				return mock.NewMockUserRepository(ctrl)
			},
			input: application.UsecaseUserRegisterInput{
				Name:     "Sr Foo",
				Email:    "foo.bar.tld",
				Password: "password",
			},
			userExpected: domain.User{},
			errExpected:  domain.NewErrInvalidEmail("foo.bar.tld"),
		},
		{
			name: "Test register user with unknown error from user repository on get by e-mail",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := mock.NewMockUserRepository(ctrl)
				repo.
					EXPECT().
					GetByEmail(gomock.Any(), gomock.Eq("foo@bar.tld")).
					Return(domain.User{}, errors.New("database gone")).
					Times(1)

				return repo
			},
			input: application.UsecaseUserRegisterInput{
				Name:     "Sr Foo",
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userExpected: domain.User{},
			errExpected:  errors.New("database gone"),
		},
		{
			name: "Test register user with user already registered by e-mail informed",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := mock.NewMockUserRepository(ctrl)
				repo.
					EXPECT().
					GetByEmail(gomock.Any(), gomock.Eq("foo@bar.tld")).
					Return(domain.User{
						ID:       "ID",
						Name:     "Sr Foo",
						Email:    "foo@bar.tld",
						Password: "password-hashed",
					}, nil).
					Times(1)

				return repo
			},
			input: application.UsecaseUserRegisterInput{
				Name:     "Sr Foo",
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userExpected: domain.User{},
			errExpected:  application.NewErrUserAlreadyRegistered("foo@bar.tld"),
		},
		{
			name: "Test register user with unknown error from user repository on persist",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := mock.NewMockUserRepository(ctrl)
				repo.
					EXPECT().
					GetByEmail(gomock.Any(), gomock.Eq("foo@bar.tld")).
					Return(domain.User{}, nil).
					Times(1)
				repo.
					EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
					Return(domain.User{}, errors.New("database gone")).
					Times(1)

				return repo
			},
			input: application.UsecaseUserRegisterInput{
				Name:     "Sr Foo",
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userExpected: domain.User{},
			errExpected:  errors.New("database gone"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.TODO()
			ctrl := gomock.NewController(t)

			usecase := application.NewUsecaseUserRegister(test.userRepositoryDependency(ctrl))
			userGot, errGot := usecase.Execute(ctx, test.input)

			assert.Equal(t, test.userExpected, userGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}
