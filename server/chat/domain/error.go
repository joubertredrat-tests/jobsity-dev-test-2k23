package domain

import "fmt"

type ErrInvalidEmail struct {
	email string
}

func NewErrInvalidEmail(email string) ErrInvalidEmail {
	return ErrInvalidEmail{
		email: email,
	}
}

func (e ErrInvalidEmail) Error() string {
	return fmt.Sprintf("Invalid e-mail got [ %s ]", e.email)
}

type ErrInvalidPasswordLength struct {
	expected uint
	got      uint
}

func NewErrInvalidPasswordLength(expected, got uint) ErrInvalidPasswordLength {
	return ErrInvalidPasswordLength{
		expected: expected,
		got:      got,
	}
}

func (e ErrInvalidPasswordLength) Error() string {
	return fmt.Sprintf("Invalid password length, expected [ %d ], got [ %d ]", e.expected, e.got)
}
