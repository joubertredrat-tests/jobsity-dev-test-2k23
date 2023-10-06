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
