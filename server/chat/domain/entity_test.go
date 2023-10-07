package domain_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	tests := []struct {
		name          string
		nameInput     string
		emailInput    string
		passwordInput string
		userExpected  domain.User
		errExpected   error
	}{
		{
			name:          "Test create user with success",
			nameInput:     "Mr Foo",
			emailInput:    "foo@bar.tld",
			passwordInput: "password",
			userExpected: domain.User{
				ID:       "",
				Name:     "Mr Foo",
				Email:    "foo@bar.tld",
				Password: "password",
			},
			errExpected: nil,
		},
		{
			name:          "Test create user with invalid e-mail",
			nameInput:     "Mr Foo",
			emailInput:    "foo.bar.tld",
			passwordInput: "password",
			userExpected:  domain.User{},
			errExpected:   domain.NewErrInvalidEmail("foo.bar.tld"),
		},
		{
			name:          "Test create user with invalid password",
			nameInput:     "Mr Foo",
			emailInput:    "foo@bar.tld",
			passwordInput: "psswd",
			userExpected:  domain.User{},
			errExpected:   domain.NewErrInvalidPasswordLength(uint(domain.PASSWORD_MIN_LENGTH), uint(5)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userGot, errGot := domain.NewUser("", test.nameInput, test.emailInput, test.passwordInput)

			assert.Equal(t, test.userExpected, userGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}

func TestMessage(t *testing.T) {
	tests := []struct {
		name            string
		userNameInput   string
		userEmailInput  string
		textInput       string
		messageExpected domain.Message
		errExpected     error
	}{
		{
			name:           "Test create message with success",
			userNameInput:  "Mr Foo",
			userEmailInput: "foo@bar.tld",
			textInput:      "I like cookies",
			messageExpected: domain.Message{
				UserName:  "Mr Foo",
				UserEmail: "foo@bar.tld",
				Text:      "I like cookies",
			},
			errExpected: nil,
		},
		{
			name:            "Test create message with invalid e-mail",
			userNameInput:   "Mr Foo",
			userEmailInput:  "foo.bar.tld",
			textInput:       "I like cookies",
			messageExpected: domain.Message{},
			errExpected:     domain.NewErrInvalidEmail("foo.bar.tld"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			messageGot, errGot := domain.NewMessage("", test.userNameInput, test.userEmailInput, test.textInput)

			assert.Equal(t, test.messageExpected, messageGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}
