package domain

import "net/mail"

type User struct {
	Name     string
	Email    string
	Password string
}

func NewUser(name, email, password string) (User, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return User{}, NewErrInvalidEmail(email)
	}

	return User{
		Name:     name,
		Email:    email,
		Password: password,
	}, nil
}
