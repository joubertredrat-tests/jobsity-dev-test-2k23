package application

type UsecaseUserRegisterInput struct {
	Name     string
	Email    string
	Password string
}

type UsecaseUserLoginInput struct {
	Email    string
	Password string
}
