package domain_test

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/internal/domain"
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
