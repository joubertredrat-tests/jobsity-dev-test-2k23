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

type ErrUserNotFoundByEmail struct {
	email string
}

func NewErrUserNotFoundByEmail(email string) ErrUserNotFoundByEmail {
	return ErrUserNotFoundByEmail{
		email: email,
	}
}

func (e ErrUserNotFoundByEmail) Error() string {
	return fmt.Sprintf("User not found by e-mail [ %s ]", e.email)
}

type ErrUserNotAuthenticated struct {
	email string
}

func NewErrUserNotAuthenticated(email string) ErrUserNotAuthenticated {
	return ErrUserNotAuthenticated{
		email: email,
	}
}

func (e ErrUserNotAuthenticated) Error() string {
	return fmt.Sprintf("User not authenticated by e-mail [ %s ]", e.email)
}
