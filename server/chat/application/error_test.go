package application_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/application"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrUserAlreadyRegistered(t *testing.T) {
	email := "foo@bar.tld"
	errExpected := "User already registered with e-mail [ foo@bar.tld ]"
	errGot := application.NewErrUserAlreadyRegistered(email)

	assert.Equal(t, errExpected, errGot.Error())
}
