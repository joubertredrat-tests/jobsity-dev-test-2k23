package application

import "fmt"

type ErrUserAlreadyRegistered struct {
	email string
}

func NewErrUserAlreadyRegistered(email string) ErrUserAlreadyRegistered {
	return ErrUserAlreadyRegistered{
		email: email,
	}
}

func (e ErrUserAlreadyRegistered) Error() string {
	return fmt.Sprintf("User already registered with e-mail [ %s ]", e.email)
}
