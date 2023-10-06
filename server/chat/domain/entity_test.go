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
			nameInput:     "foo",
			emailInput:    "foo@bar.tld",
			passwordInput: "password",
			userExpected: domain.User{
				ID:       "",
				Name:     "foo",
				Email:    "foo@bar.tld",
				Password: "password",
			},
			errExpected: nil,
		},
		{
			name:          "Test create user with invalid e-mail",
			nameInput:     "foo",
			emailInput:    "foo.bar.tld",
			passwordInput: "password",
			userExpected:  domain.User{},
			errExpected:   domain.NewErrInvalidEmail("foo.bar.tld"),
		},
		{
			name:          "Test create user with invalid password",
			nameInput:     "foo",
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
