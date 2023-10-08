package domain_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrInvalidEmail(t *testing.T) {
	email := "foo.bar.tld"
	errExpected := "Invalid e-mail got [ foo.bar.tld ]"
	errGot := domain.NewErrInvalidEmail(email)

	assert.Equal(t, errExpected, errGot.Error())
}

func TestErrInvalidPasswordLength(t *testing.T) {
	expected := uint(8)
	got := uint(6)
	errExpected := "Invalid password length, expected [ 8 ], got [ 6 ]"
	errGot := domain.NewErrInvalidPasswordLength(expected, got)

	assert.Equal(t, errExpected, errGot.Error())
}

func TestErrUserNotFoundByEmail(t *testing.T) {
	email := "foo@bar.tld"
	errExpected := "User not found by e-mail [ foo@bar.tld ]"
	errGot := domain.NewErrUserNotFoundByEmail(email)

	assert.Equal(t, errExpected, errGot.Error())
}

func TestErrUserNotAuthenticated(t *testing.T) {
	email := "foo@bar.tld"
	errExpected := "User not authenticated by e-mail [ foo@bar.tld ]"
	errGot := domain.NewErrUserNotAuthenticated(email)

	assert.Equal(t, errExpected, errGot.Error())
}

func TestErrPaginationPage(t *testing.T) {
	errExpected := "Invalid pagination page [ 0 ]"
	errGot := domain.NewErrPaginationPage(0)

	assert.Equal(t, errExpected, errGot.Error())
}

func TestErrPaginationItemsPerPage(t *testing.T) {
	errExpected := "Invalid pagination items per page, expected between [ 10 ] and [ 100 ], got [ 107 ]"
	errGot := domain.NewErrPaginationItemsPerPage(10, 100, 107)

	assert.Equal(t, errExpected, errGot.Error())
}
