package application_test

import (
	"context"
	"errors"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/application"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"joubertredrat-tests/jobsity-dev-test-2k23/pkg"
	"joubertredrat-tests/jobsity-dev-test-2k23/pkg/chat/mock"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var errDatabaseGone = errors.New("database gone")
var errTokenGeneration = errors.New("fail on token generation")

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
					Return(domain.User{}, errDatabaseGone).
					Times(1)

				return repo
			},
			input: application.UsecaseUserRegisterInput{
				Name:     "Sr Foo",
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userExpected: domain.User{},
			errExpected:  errDatabaseGone,
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
			name: "Test register user with unknown error from user repository on create",
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
					Return(domain.User{}, errDatabaseGone).
					Times(1)

				return repo
			},
			input: application.UsecaseUserRegisterInput{
				Name:     "Sr Foo",
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userExpected: domain.User{},
			errExpected:  errDatabaseGone,
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

func TestUsecaseUserLogin(t *testing.T) {
	tests := []struct {
		name                     string
		userRepositoryDependency func(ctrl *gomock.Controller) domain.UserRepository
		tokenServiceDependency   func(ctrl *gomock.Controller) domain.TokenService
		input                    application.UsecaseUserLoginInput
		userTokenExpected        domain.UserToken
		errExpected              error
	}{
		{
			name: "Test user login with success",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := mock.NewMockUserRepository(ctrl)
				repo.
					EXPECT().
					GetAuthenticated(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
					Return(domain.User{
						ID:       "ID",
						Name:     "Sr Foo",
						Email:    "foo@bar.tld",
						Password: "password-hashed",
					}, nil).
					Times(1)

				return repo
			},
			tokenServiceDependency: func(ctrl *gomock.Controller) domain.TokenService {
				service := mock.NewMockTokenService(ctrl)
				service.
					EXPECT().
					Generate(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
					Return(domain.UserToken{
						AccessToken: "token",
					}, nil).
					Times(1)

				return service
			},
			input: application.UsecaseUserLoginInput{
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userTokenExpected: domain.UserToken{
				AccessToken: "token",
			},
			errExpected: nil,
		},
		{
			name: "Test user login with invalid e-mail",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				return mock.NewMockUserRepository(ctrl)
			},
			tokenServiceDependency: func(ctrl *gomock.Controller) domain.TokenService {
				return mock.NewMockTokenService(ctrl)
			},
			input: application.UsecaseUserLoginInput{
				Email:    "foo.bar.tld",
				Password: "password",
			},
			userTokenExpected: domain.UserToken{},
			errExpected:       domain.NewErrInvalidEmail("foo.bar.tld"),
		},
		{
			name: "Test user login with unknown error from user repository on get authenticated",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := mock.NewMockUserRepository(ctrl)
				repo.
					EXPECT().
					GetAuthenticated(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
					Return(domain.User{}, errDatabaseGone).
					Times(1)

				return repo
			},
			tokenServiceDependency: func(ctrl *gomock.Controller) domain.TokenService {
				return mock.NewMockTokenService(ctrl)
			},
			input: application.UsecaseUserLoginInput{
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userTokenExpected: domain.UserToken{},
			errExpected:       errDatabaseGone,
		},
		{
			name: "Test user login with user not found by e-mail error from user repository on get authenticated",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := mock.NewMockUserRepository(ctrl)
				repo.
					EXPECT().
					GetAuthenticated(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
					Return(domain.User{}, domain.NewErrUserNotFoundByEmail("foo@bar.tld")).
					Times(1)

				return repo
			},
			tokenServiceDependency: func(ctrl *gomock.Controller) domain.TokenService {
				return mock.NewMockTokenService(ctrl)
			},
			input: application.UsecaseUserLoginInput{
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userTokenExpected: domain.UserToken{},
			errExpected:       domain.NewErrUserNotFoundByEmail("foo@bar.tld"),
		},
		{
			name: "Test user login with unknown error from token service on generate",
			userRepositoryDependency: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := mock.NewMockUserRepository(ctrl)
				repo.
					EXPECT().
					GetAuthenticated(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
					Return(domain.User{
						ID:       "ID",
						Name:     "Sr Foo",
						Email:    "foo@bar.tld",
						Password: "password-hashed",
					}, nil).
					Times(1)

				return repo
			},
			tokenServiceDependency: func(ctrl *gomock.Controller) domain.TokenService {
				service := mock.NewMockTokenService(ctrl)
				service.
					EXPECT().
					Generate(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).
					Return(domain.UserToken{}, errTokenGeneration).
					Times(1)

				return service
			},
			input: application.UsecaseUserLoginInput{
				Email:    "foo@bar.tld",
				Password: "password",
			},
			userTokenExpected: domain.UserToken{},
			errExpected:       errTokenGeneration,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.TODO()
			ctrl := gomock.NewController(t)

			usecase := application.NewUsecaseUserLogin(
				test.userRepositoryDependency(ctrl),
				test.tokenServiceDependency(ctrl),
			)
			userTokenGot, errGot := usecase.Execute(ctx, test.input)

			assert.Equal(t, test.userTokenExpected, userTokenGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}

func TestUsecaseMessageCreate(t *testing.T) {
	tests := []struct {
		name                        string
		messageRepositoryDependency func(ctrl *gomock.Controller) domain.MessageRepository
		messageEventDependency      func(ctrl *gomock.Controller) domain.MessageEvent
		input                       application.UsecaseMessageCreateInput
		messageExpected             domain.Message
		errExpected                 error
	}{
		{
			name: "Test create normal message with success",
			messageRepositoryDependency: func(ctrl *gomock.Controller) domain.MessageRepository {
				repo := mock.NewMockMessageRepository(ctrl)

				datetime := "2023-10-07 14:27:51"
				time, _ := pkg.TimeFromCanonical(&datetime)

				repo.
					EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Message{})).
					Return(domain.Message{
						ID:        "ID",
						UserName:  "Sr Foo",
						UserEmail: "foo@bar.tld",
						Text:      "I like cookies",
						Datetime:  *time,
					}, nil).
					Times(1)

				return repo
			},
			messageEventDependency: func(ctrl *gomock.Controller) domain.MessageEvent {
				return mock.NewMockMessageEvent(ctrl)
			},
			input: application.UsecaseMessageCreateInput{
				UserName:    "Sr Foo",
				UserEmail:   "foo@bar.tld",
				MessageText: "I like cookies",
			},
			messageExpected: func() domain.Message {
				datetime := "2023-10-07 14:27:51"
				time, _ := pkg.TimeFromCanonical(&datetime)

				return domain.Message{
					ID:        "ID",
					UserName:  "Sr Foo",
					UserEmail: "foo@bar.tld",
					Text:      "I like cookies",
					Datetime:  *time,
				}
			}(),
			errExpected: nil,
		},
		{
			name: "Test receive command message with success",
			messageRepositoryDependency: func(ctrl *gomock.Controller) domain.MessageRepository {
				return mock.NewMockMessageRepository(ctrl)
			},
			messageEventDependency: func(ctrl *gomock.Controller) domain.MessageEvent {
				event := mock.NewMockMessageEvent(ctrl)

				event.
					EXPECT().
					StockCommandReceived(gomock.Any(), gomock.AssignableToTypeOf(domain.Message{})).
					Return(nil).
					Times(1)

				return event
			},
			input: application.UsecaseMessageCreateInput{
				UserName:    "Sr Foo",
				UserEmail:   "foo@bar.tld",
				MessageText: "/stock=APPL.US",
			},
			messageExpected: domain.Message{
				UserName:  "Sr Foo",
				UserEmail: "foo@bar.tld",
				Text:      "/stock=APPL.US",
			},
			errExpected: nil,
		},
		{
			name: "Test create message with invalid e-mail",
			messageRepositoryDependency: func(ctrl *gomock.Controller) domain.MessageRepository {
				return mock.NewMockMessageRepository(ctrl)
			},
			messageEventDependency: func(ctrl *gomock.Controller) domain.MessageEvent {
				return mock.NewMockMessageEvent(ctrl)
			},
			input: application.UsecaseMessageCreateInput{
				UserName:    "Sr Foo",
				UserEmail:   "foo.bar.tld",
				MessageText: "I like cookies",
			},
			messageExpected: domain.Message{},
			errExpected:     domain.NewErrInvalidEmail("foo.bar.tld"),
		},
		{
			name: "Test create message with unknown error from message repository on create",
			messageRepositoryDependency: func(ctrl *gomock.Controller) domain.MessageRepository {
				repo := mock.NewMockMessageRepository(ctrl)

				repo.
					EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Message{})).
					Return(domain.Message{}, errDatabaseGone).
					Times(1)

				return repo
			},
			messageEventDependency: func(ctrl *gomock.Controller) domain.MessageEvent {
				return mock.NewMockMessageEvent(ctrl)
			},
			input: application.UsecaseMessageCreateInput{
				UserName:    "Sr Foo",
				UserEmail:   "foo@bar.tld",
				MessageText: "I like cookies",
			},
			messageExpected: domain.Message{},
			errExpected:     errDatabaseGone,
		},
		{
			name: "Test receive command message with unknown error from message event on stock command received",
			messageRepositoryDependency: func(ctrl *gomock.Controller) domain.MessageRepository {
				return mock.NewMockMessageRepository(ctrl)
			},
			messageEventDependency: func(ctrl *gomock.Controller) domain.MessageEvent {
				event := mock.NewMockMessageEvent(ctrl)

				event.
					EXPECT().
					StockCommandReceived(gomock.Any(), gomock.AssignableToTypeOf(domain.Message{})).
					Return(errors.New("message broker is not connected")).
					Times(1)

				return event
			},
			input: application.UsecaseMessageCreateInput{
				UserName:    "Sr Foo",
				UserEmail:   "foo@bar.tld",
				MessageText: "/stock=APPL.US",
			},
			messageExpected: domain.Message{},
			errExpected:     errors.New("message broker is not connected"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.TODO()
			ctrl := gomock.NewController(t)

			usecase := application.NewUsecaseMessageCreate(
				test.messageRepositoryDependency(ctrl),
				test.messageEventDependency(ctrl),
			)
			messageGot, errGot := usecase.Execute(ctx, test.input)

			assert.Equal(t, test.messageExpected, messageGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}

func TestUsecaseMessageList(t *testing.T) {
	tests := []struct {
		name                        string
		messageRepositoryDependency func(ctrl *gomock.Controller) domain.MessageRepository
		input                       application.UsecaseMessagesListInput
		messagesExpected            []domain.Message
		errExpected                 error
	}{
		{
			name: "Test list messages with success",
			messageRepositoryDependency: func(ctrl *gomock.Controller) domain.MessageRepository {
				repo := mock.NewMockMessageRepository(ctrl)

				datetime1 := "2023-10-07 14:27:51"
				time1, _ := pkg.TimeFromCanonical(&datetime1)

				datetime2 := "2023-10-07 12:57:06"
				time2, _ := pkg.TimeFromCanonical(&datetime2)

				repo.
					EXPECT().
					List(
						gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
						gomock.AssignableToTypeOf(domain.Pagination{}),
					).
					Return([]domain.Message{
						{
							ID:        "ID2",
							UserName:  "Mr Qux",
							UserEmail: "qux@bar.tld",
							Text:      "Me too",
							Datetime:  *time2,
						},
						{
							ID:        "ID1",
							UserName:  "Sr Foo",
							UserEmail: "foo@bar.tld",
							Text:      "I like cookies",
							Datetime:  *time1,
						},
					}, nil).
					Times(1)

				return repo
			},
			input: application.UsecaseMessagesListInput{
				Page:         1,
				ItemsPerPage: 50,
			},
			messagesExpected: func() []domain.Message {
				datetime1 := "2023-10-07 14:27:51"
				time1, _ := pkg.TimeFromCanonical(&datetime1)

				datetime2 := "2023-10-07 12:57:06"
				time2, _ := pkg.TimeFromCanonical(&datetime2)

				return []domain.Message{
					{
						ID:        "ID2",
						UserName:  "Mr Qux",
						UserEmail: "qux@bar.tld",
						Text:      "Me too",
						Datetime:  *time2,
					},
					{
						ID:        "ID1",
						UserName:  "Sr Foo",
						UserEmail: "foo@bar.tld",
						Text:      "I like cookies",
						Datetime:  *time1,
					},
				}
			}(),
			errExpected: nil,
		},
		{
			name: "Test list messages with invalid page",
			messageRepositoryDependency: func(ctrl *gomock.Controller) domain.MessageRepository {
				return mock.NewMockMessageRepository(ctrl)
			},
			input: application.UsecaseMessagesListInput{
				Page:         0,
				ItemsPerPage: 50,
			},
			messagesExpected: []domain.Message{},
			errExpected:      domain.NewErrPaginationPage(0),
		},
		{
			name: "Test list messages with unknown error from message repository on list",
			messageRepositoryDependency: func(ctrl *gomock.Controller) domain.MessageRepository {
				repo := mock.NewMockMessageRepository(ctrl)

				repo.
					EXPECT().
					List(
						gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
						gomock.AssignableToTypeOf(domain.Pagination{}),
					).
					Return([]domain.Message{}, errDatabaseGone).
					Times(1)

				return repo
			},
			input: application.UsecaseMessagesListInput{
				Page:         1,
				ItemsPerPage: 50,
			},
			messagesExpected: []domain.Message{},
			errExpected:      errDatabaseGone,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.TODO()
			ctrl := gomock.NewController(t)

			usecase := application.NewUsecaseMessageList(test.messageRepositoryDependency(ctrl))
			messagesGot, errGot := usecase.Execute(ctx, test.input)

			assert.Equal(t, test.messagesExpected, messagesGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}
