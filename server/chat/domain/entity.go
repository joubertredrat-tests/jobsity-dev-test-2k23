package domain

import "net/mail"

const PASSWORD_MIN_LENGTH = 8

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func NewUser(id, name, email, password string) (User, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return User{}, NewErrInvalidEmail(email)
	}
	if password != "" && PASSWORD_MIN_LENGTH > len(password) {
		return User{}, NewErrInvalidPasswordLength(uint(PASSWORD_MIN_LENGTH), uint(len(password)))
	}

	return User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}, nil
}

type UserToken struct {
	AccessToken string
}

func NewUserToken(accesToken string) UserToken {
	return UserToken{
		AccessToken: accesToken,
	}
}
